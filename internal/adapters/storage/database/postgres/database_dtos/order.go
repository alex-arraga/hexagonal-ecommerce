package database_dtos

import (
	"go-ecommerce/internal/adapters/storage/database/postgres/models"
	"go-ecommerce/internal/core/domain"
)

// domain.Order -> DB model
func ConvertOrderDomainToModel(o *domain.Order) *models.OrderModel {
	return &models.OrderModel{
		ID:                o.ID,
		Providers:         o.Providers,
		UserID:            o.UserID,
		PaymentID:         o.PaymentID,
		SecureToken:       o.SecureToken,
		ExternalReference: o.ExternalReference,
		Currency:          o.Currency,
		SubTotal:          o.SubTotal,
		Disscount:         o.Disscount,
		DisscountType:     o.DisscountType,
		Total:             o.Total,
		Paid:              o.Paid,
		PayStatus:         o.PayStatus,
		PayStatusDetail:   o.PayStatusDetail,
		CreatedAt:         o.CreatedAt,
		UpdatedAt:         o.UpdatedAt,
		ExpiresAt:         o.ExpiresAt,
	}
}

// domain.Orders -> DB models
func ConvertOrdersDomainToModels(orders []*domain.Order) []*models.OrderModel {
	var ordersModels []*models.OrderModel

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
			Disscount:         o.Disscount,
			DisscountType:     o.DisscountType,
			Total:             o.Total,
			Paid:              o.Paid,
			PayStatus:         o.PayStatus,
			PayStatusDetail:   o.PayStatusDetail,
			CreatedAt:         o.CreatedAt,
			UpdatedAt:         o.UpdatedAt,
			ExpiresAt:         o.ExpiresAt,
		})
	}

	return ordersModels
}

// DB model -> domain.Order
func ConvertOrderModelToDomain(p *models.OrderModel) *domain.Order {
	return &domain.Order{
		ID:                p.ID,
		Providers:         p.Providers,
		UserID:            p.UserID,
		PaymentID:         p.PaymentID,
		SecureToken:       p.SecureToken,
		ExternalReference: p.ExternalReference,
		Currency:          p.Currency,
		SubTotal:          p.SubTotal,
		Disscount:         p.Disscount,
		DisscountType:     p.DisscountType,
		Total:             p.Total,
		Paid:              p.Paid,
		PayStatus:         p.PayStatus,
		PayStatusDetail:   p.PayStatusDetail,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		ExpiresAt:         p.ExpiresAt,
	}
}

// DB models -> domain.Orders
func ConvertOrdersModelsToDomain(orders []*models.OrderModel) []*domain.Order {
	var ordersDomain []*domain.Order

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
			Disscount:         o.Disscount,
			DisscountType:     o.DisscountType,
			Total:             o.Total,
			Paid:              o.Paid,
			PayStatus:         o.PayStatus,
			PayStatusDetail:   o.PayStatusDetail,
			CreatedAt:         o.CreatedAt,
			UpdatedAt:         o.UpdatedAt,
			ExpiresAt:         o.ExpiresAt,
		})
	}

	return ordersDomain
}
