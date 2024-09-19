package bitmex

type OrderStatusType string

const (
	OrderStatusTypeNew             OrderStatusType = "New"
	OrderStatusTypePartiallyFilled OrderStatusType = "PartiallyFilled"
	OrderStatusTypeFilled          OrderStatusType = "Filled"
)
