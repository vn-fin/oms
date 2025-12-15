package typing

type ActionType string

const (
	ActionBuy  ActionType = "B"
	ActionSell ActionType = "S"
)

func (t ActionType) Valid() bool {
	switch t {
	case ActionBuy, ActionSell:
		return true
	default:
		return false
	}
}
