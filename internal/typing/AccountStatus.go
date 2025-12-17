package typing

type AccountStatus string

const (
	StatusActive   AccountStatus = "active"
	StatusDisabled AccountStatus = "disabled"
)

func (t AccountStatus) Valid() bool {
	switch t {
	case StatusActive, StatusDisabled:
		return true
	default:
		return false
	}
}
