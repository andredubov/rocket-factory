package model

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusAssembled OrderStatus = "ASSEMBLED"
)

func NewOrderStatus(orderStatus string) (OrderStatus, error) {
	switch orderStatus {
	case "PENDING_PAYMENT":
		return OrderStatusPending, nil
	case "PAID":
		return OrderStatusPaid, nil
	case "CANCELLED":
		return OrderStatusCancelled, nil
	case "ASSEMBLED":
		return OrderStatusAssembled, nil
	default:
		return OrderStatusPending, ErrInvalidOrderStatus
	}
}

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

func NewPaymentMethod(paymentMethod string) (PaymentMethod, error) {
	switch paymentMethod {
	case "UNKNOWN":
		return PaymentMethodUnknown, nil
	case "CARD":
		return PaymentMethodCard, nil
	case "SBP":
		return PaymentMethodSBP, nil
	case "CREDIT_CARD":
		return PaymentMethodCreditCard, nil
	case "INVESTOR_MONEY":
		return PaymentMethodInvestorMoney, nil
	default:
		return PaymentMethodUnknown, ErrInvalidPaymentMethod
	}
}
