package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.Order -> DB model
func ConvertOrderDomainToModel(o *domain.Order) *models.OrderModel {
	items := make([]models.OrderProductModel, len(o.Items))
	for i, item := range o.Items {
		items[i] = models.OrderProductModel{
			ID:        item.ID,
			OrderID:   o.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return &models.OrderModel{
		ID:                o.ID,
		Providers:         o.Providers,
		UserID:            o.UserID,
		PaymentID:         o.PaymentID,
		SecureToken:       o.SecureToken,
		ExternalReference: o.ExternalReference,
		Currency:          o.Currency,
		SubTotal:          o.SubTotal,
		Discount:          o.Discount,
		Total:             o.Total,
		Paid:              o.Paid,
		Fee:               *o.Fee,
		Installments:      *o.Installments,
		PayMethod:         o.PayMethod,
		PayResource:       o.PayResource,
		NetReceivedAmount: *o.NetReceivedAmount,
		PayStatus:         o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
		ExpiresAt:         o.ExpiresAt,
		PaidAt:            o.PaidAt,
		Items:             items,
	}
}

// domain.Orders -> DB models
func ConvertOrdersDomainToModels(orders []*domain.Order) []*models.OrderModel {
	var ordersModels []*models.OrderModel

	var items []models.OrderProductModel
	for _, o := range orders {
		items = make([]models.OrderProductModel, len(o.Items))
		for i, item := range o.Items {
			items[i] = models.OrderProductModel{
				ID:        item.ID,
				OrderID:   o.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}
	}

	for _, o := range orders {
		ordersModels = append(ordersModels, &models.OrderModel{
			ID:                o.ID,
			Providers:         o.Providers,
			UserID:            o.UserID,
			PaymentID:         o.PaymentID,
			SecureToken:       o.SecureToken,
			ExternalReference: o.ExternalReference,
			Currency:          o.Currency,
			SubTotal:          o.SubTotal,
			Discount:          o.Discount,
			Total:             o.Total,
			Paid:              o.Paid,
			Fee:               *o.Fee,
			Installments:      *o.Installments,
			PayMethod:         o.PayMethod,
			PayResource:       o.PayResource,
			NetReceivedAmount: *o.NetReceivedAmount,
			PayStatus:         o.PayStatus,
			PayStatusDetail:   o.PayStatusDetail,
			CreatedAt:         o.CreatedAt,
			UpdatedAt:         o.UpdatedAt,
			ExpiresAt:         o.ExpiresAt,
			PaidAt:            o.PaidAt,
			Items:             items,
		})
	}

	return ordersModels
}

// DB model -> domain.Order
func ConvertOrderModelToDomain(o *models.OrderModel) *domain.Order {
	items := make([]domain.OrderProduct, len(o.Items))
	for i, item := range o.Items {
		items[i] = domain.OrderProduct{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	return &domain.Order{
		ID:                o.ID,
		Providers:         o.Providers,
		UserID:            o.UserID,
		PaymentID:         o.PaymentID,
		SecureToken:       o.SecureToken,
		ExternalReference: o.ExternalReference,
		Currency:          o.Currency,
		SubTotal:          o.SubTotal,
		Discount:          o.Discount,
		Total:             o.Total,
		Paid:              o.Paid,
		Fee:               &o.Fee,
		Installments:      &o.Installments,
		PayMethod:         o.PayMethod,
		PayResource:       o.PayResource,
		NetReceivedAmount: &o.NetReceivedAmount,
		PayStatus:         o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
		ExpiresAt:         o.ExpiresAt,
		PaidAt:            o.PaidAt,
		Items:             items,
	}
}

// DB models -> domain.Orders
func ConvertOrdersModelsToDomain(orders []*models.OrderModel) []*domain.Order {
	var ordersDomain []*domain.Order
	items := make([]domain.OrderProduct, 0)

	for _, orders := range orders {
		items = make([]domain.OrderProduct, len(orders.Items))
		for i, item := range orders.Items {
			items[i] = domain.OrderProduct{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}
	}

	for _, o := range orders {
		ordersDomain = append(ordersDomain, &domain.Order{
			ID:                o.ID,
			Providers:         o.Providers,
			UserID:            o.UserID,
			PaymentID:         o.PaymentID,
			SecureToken:       o.SecureToken,
			ExternalReference: o.ExternalReference,
			Currency:          o.Currency,
			SubTotal:          o.SubTotal,
			Discount:          o.Discount,
			Total:             o.Total,
			Paid:              o.Paid,
			Fee:               &o.Fee,
			Installments:      &o.Installments,
			PayMethod:         o.PayMethod,
			PayResource:       o.PayResource,
			NetReceivedAmount: &o.NetReceivedAmount,
			PayStatus:         o.PayStatus,
			PayStatusDetail:   o.PayStatusDetail,
			CreatedAt:         o.CreatedAt,
			UpdatedAt:         o.UpdatedAt,
			ExpiresAt:         o.ExpiresAt,
			PaidAt:            o.PaidAt,
			Items:             items,
		})
	}

	return ordersDomain
}
