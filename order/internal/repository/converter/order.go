package converter

import (
	"github.com/andredubov/rocket-factory/order/internal/model"
	repoModel "github.com/andredubov/rocket-factory/order/internal/repository/model"
)

// OrderToRepoModel converts a domain model (repo.Order) to a repository model (repoModel.Order).
func OrderToRepoModel(order model.Order) repoModel.Order {
	var paymentInfo *repoModel.PaymentInfo
	if order.PaymentInfo != nil {
		paymentInfo = &repoModel.PaymentInfo{
			TransactionUUID: order.PaymentInfo.TransactionUUID,
			PaymentMethod:   repoModel.PaymentMethod(order.PaymentInfo.PaymentMethod),
		}
	}

	return repoModel.Order{
		OrderUUID:   order.OrderUUID,
		UserUUID:    order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  order.TotalPrice,
		PaymentInfo: paymentInfo,
		Status:      repoModel.OrderStatus(order.Status),
	}
}

// OrderToModel converts a repository model (repoModel.Order) to a domain model (model.Order).
func OrderToModel(order repoModel.Order) *model.Order {
	var paymentInfo *model.PaymentInfo
	if order.PaymentInfo != nil {
		paymentInfo = &model.PaymentInfo{
			TransactionUUID: order.PaymentInfo.TransactionUUID,
			PaymentMethod:   model.PaymentMethod(order.PaymentInfo.PaymentMethod),
		}
	}

	return &model.Order{
		OrderUUID:   order.OrderUUID,
		UserUUID:    order.UserUUID,
		PartUUIDs:   order.PartUUIDs,
		TotalPrice:  order.TotalPrice,
		PaymentInfo: paymentInfo,
		Status:      model.OrderStatus(order.Status),
	}
}

// OrdersToModel converts a slice of repository order models (repoModel.Order) to a slice of domain order models (model.Order).
func OrdersToModel(orders []repoModel.Order) []model.Order {
	result := make([]model.Order, 0, len(orders))

	for _, order := range orders {
		domainOrder := OrderToModel(order)
		if domainOrder != nil {
			result = append(result, *domainOrder)
		}
	}

	return result
}
