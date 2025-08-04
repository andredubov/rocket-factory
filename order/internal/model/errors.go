package model

import "errors"

var (
	ErrInvalidOrderStatus    = errors.New("invalid order status")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrOrderAlreadyExists    = errors.New("order already exists")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrOrderHasNoParts       = errors.New("order has no parts")
	ErrInvalidPartFilter     = errors.New("invalid part filter")
)
