package model

type OrderPaidEvent struct {
	UUID            string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   PaymentMethod
	TransactionUUID string
}

type OrderAssembledEvent struct {
	UUID         string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
