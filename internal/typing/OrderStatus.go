package typing

type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusCanceled      OrderStatus = "canceled"
	OrderStatusPartialFilled OrderStatus = "partial-filled"
)

func (s OrderStatus) Valid() bool {
	switch s {
	case OrderStatusPending, OrderStatusCanceled, OrderStatusPartialFilled:
		return true
	default:
		return false
	}
}
