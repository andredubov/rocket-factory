package model

type PaymentMethod int32

// Valid PaymentMethod values
const (
	PaymentMethodUnknown       PaymentMethod = 0
	PaymentMethodCard          PaymentMethod = 1
	PaymentMethodSbp           PaymentMethod = 2
	PaymentMethodreditCard     PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

// IsValid checks if the PaymentMethod has a valid value
func (os PaymentMethod) IsValid() bool {
	switch os {
	case PaymentMethodUnknown, PaymentMethodCard, PaymentMethodSbp, PaymentMethodreditCard, PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}

// Payment represents a payment transaction in the domain layer.
type Payment struct {
	UserUuid      string
	OrderUuid     string
	PaymentMethod PaymentMethod
}
