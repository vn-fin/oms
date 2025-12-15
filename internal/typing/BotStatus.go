package typing

type BotStatus string

const BotStatusActive BotStatus = "active"
const BotStatusDisabled BotStatus = "disabled"

func (e BotStatus) Valid() bool {
	switch e {
	case BotStatusActive, BotStatusDisabled:
		return true
	default:
		return false
	}
}

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
