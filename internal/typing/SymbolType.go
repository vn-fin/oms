package typing

type SymbolType string

const (
	VnStock  SymbolType = "VnStock"
	VnFuture SymbolType = "VnFuture"
)

func (t SymbolType) Valid() bool {
	switch t {
	case VnStock, VnFuture:
		return true
	default:
		return false
	}
}
