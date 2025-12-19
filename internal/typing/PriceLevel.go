package typing

type PriceLevel string

const (
	PriceLevelMid   PriceLevel = "mid"
	PriceLevelAsk01 PriceLevel = "ask01"
	PriceLevelAsk02 PriceLevel = "ask02"
	PriceLevelAsk03 PriceLevel = "ask03"
	PriceLevelBid01 PriceLevel = "bid01"
	PriceLevelBid02 PriceLevel = "bid02"
	PriceLevelBid03 PriceLevel = "bid03"
	PriceLevelCeil  PriceLevel = "ceil"
	PriceLevelFloor PriceLevel = "floor"
)

func (s PriceLevel) Valid() bool {
	switch s {
	case PriceLevelMid, PriceLevelAsk01, PriceLevelAsk02, PriceLevelAsk03, PriceLevelBid01, PriceLevelBid02, PriceLevelBid03, PriceLevelCeil, PriceLevelFloor:
		return true
	default:
		return false
	}
}
