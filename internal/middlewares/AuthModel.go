package middlewares

type AuthModel struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CashLimit int    `json:"cash_limit"`
}
