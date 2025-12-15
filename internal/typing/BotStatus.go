package typing

type RecordStatus string

const RecordStatusEnabled RecordStatus = "enabled"
const RecordStatusDisabled RecordStatus = "disabled"
const RecordStatusRemoved RecordStatus = "removed"

func (e RecordStatus) Valid() bool {
	switch e {
	case RecordStatusEnabled, RecordStatusDisabled, RecordStatusRemoved:
		return true
	default:
		return false
	}
}
