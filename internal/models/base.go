package models

// Pagination represents the pagination information
type Pagination struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

// DefaultResponseModel represents the default response model
type DefaultResponseModel struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"` // Use interface{} to allow any type for data
	// The pagination information
	Page *Pagination `json:"page"`
}
