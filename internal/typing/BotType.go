package typing

type BotType string

const BotTypeIntraday BotType = "intraday"
const BotTypeDaily BotType = "daily"
const BotTypeHFT BotType = "hft"

func (e BotType) Valid() bool {
	switch e {
	case BotTypeIntraday, BotTypeDaily, BotTypeHFT:
		return true
	default:
		return false
	}
}
