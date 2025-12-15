package typing

type UserRole string

const (
	UserRoleAdmin  UserRole = "admin"
	UserRolePM     UserRole = "pm"
	UserRoleTrader UserRole = "trader"
)

func (u UserRole) Valid() bool {
	switch u {
	case UserRoleAdmin, UserRoleTrader, UserRolePM:
		return true
	default:
		return false
	}
}
