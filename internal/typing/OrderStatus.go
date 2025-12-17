package typing

type OrderStatus string

const (
	OrderStatusCreated       OrderStatus = "created"
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusFilled        OrderStatus = "filled"
	OrderStatusPartialFilled OrderStatus = "partial-filled"
	OrderStatusCanceled      OrderStatus = "canceled"
)

func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusCreated, OrderStatusPending, OrderStatusFilled, OrderStatusPartialFilled, OrderStatusCanceled:
		return true
	default:
		return false
	}
}
