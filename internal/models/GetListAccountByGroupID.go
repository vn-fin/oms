package models

type GetListCredentialByGroup struct {
	GroupID     string   `json:"group_id"`
	UserID      string   `json:"user_id"`
	Credentials []string `json:"credential_id"`
}
