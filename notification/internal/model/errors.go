package model

import (
	"errors"
)

var (
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
