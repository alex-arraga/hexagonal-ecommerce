package services

import (
	"context"
	"errors"
	"fmt"
	"go-ecommerce/internal/adapters/mercadopago/mp_dtos"
	"go-ecommerce/internal/core/domain"
	"go-ecommerce/internal/core/ports"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	userRepo    ports.UserRepository
	orderRepo   ports.OrderRepository
	productRepo ports.ProductRepository
	mp          ports.PaymentProvider
}

func NewPaymentService(userRepo ports.UserRepository, orderRepo ports.OrderRepository, productRepo ports.ProductRepository, mp ports.PaymentProvider) ports.PaymentService {
	return &PaymentService{
		userRepo:    userRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		mp:          mp,
	}
}

// CreatePayment implements ports.PaymentProvider.
func (p *PaymentService) StartPayment(ctx context.Context, orderId uuid.UUID) (*string, error) {
	order, err := p.orderRepo.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.GetUserByID(ctx, order.UserID)
	if err != nil {
		return nil, err
	}

	// search products of order and generate mp items
	items := make([]mp_dtos.MpItem, 0)

	for _, orderItem := range order.Items {
		product, err := p.productRepo.GetProductById(ctx, orderItem.ProductID)
		if err != nil {
			return nil, err
		}

		items = append(items, mp_dtos.MpItem{
			ID:          orderItem.ProductID.String(),
			Title:       product.Name,
			Description: product.SKU,
			CategoryID:  fmt.Sprint(product.CategoryID),
			CurrencyID:  fmt.Sprint(order.Currency),
			Quantity:    int(orderItem.Quantity),
			UnitPrice:   product.Price,
		})
	}

	// generate mercado pago preference
	preferece := p.mp.GeneratePreference(ctx, order, items, user)

	redirectUrl, err := p.mp.GenerateNewPayment(ctx, preferece)
	if err != nil {
		return nil, err
	}
	return redirectUrl, nil
}

// VerifyPayment implements ports.PaymentProvider.
func (p *PaymentService) VerifyPayment(ctx context.Context, id *string, topic *string) error {
	// tiene que retornar el id y a partir de ahi actualizar la orden
	payment, err := p.mp.VerifyPayment(ctx, id, topic)
	if err != nil {
		return err
	}

	if payment.ExternalReference == "" {
		return errors.New("external reference missing in payment")
	}

	// parse external reference (order.ID generated in mp preference) to uuid
	parsedExtRef, err := uuid.Parse(payment.ExternalReference)
	if err != nil {
		return err
	}

	// Get order to updated with payment data
	order, err := p.orderRepo.GetOrderById(ctx, parsedExtRef)
	if err != nil {
		return err
	}

	// if order has been updated, return
	if order.PaymentID != nil {
		return nil
	}

	// validations and handling errors
	// external_reference must be equal than order id
	if parsedExtRef != order.ID {
		return fmt.Errorf("payment: %s external_reference does not match with order: %s", payment.ExternalReference, order.ID)
	}

	// avoids updating a order with an approved payment but that was never was credited due to account errors or holds
	if payment.TransactionDetails.NetReceivedAmount <= 0 {
		return fmt.Errorf("net received amount is 0 or less for payment: %v", payment.ID)
	}

	// if the order has already been updated with this payment_id, return and do nothing
	strPaymentId := strconv.Itoa(payment.ID)
	if order.PaymentID == &strPaymentId {
		return fmt.Errorf("payment: %v already processed for order: %s", payment.ID, order.ID)
	}

	// update order with payment data
	successfullyPayment := payment.Status == domain.Approved && payment.StatusDetail == domain.Accredited
	fee := payment.TransactionAmount - payment.TransactionDetails.NetReceivedAmount

	var paidAt time.Time
	if successfullyPayment {
		paidAt = time.Now()
	}

	var isPaid bool
	if successfullyPayment {
		isPaid = true
	}

	dataToUpdate := domain.UpdateOrderInputs{
		PayStatus:         domain.PayStatus(payment.Status),
		PayStatusDetail:   domain.PayStatusDetail(payment.StatusDetail),
		PaymentID:         fmt.Sprint(payment.ID),
		PayMethod:         payment.PayMethod.ID,
		PayResource:       payment.PayMethod.Type,
		Installments:      payment.Installments,
		ExternalReference: payment.ExternalReference,
		NetReceivedAmount: payment.TransactionDetails.NetReceivedAmount,
		Fee:               fee,
		PaidAt:            &paidAt,
		Paid:              isPaid,
	}

	err = order.UpdateOrder(dataToUpdate)
	if err != nil {
		return err
	}

	_, err = p.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return err
	}
	return nil
}
