package pagination

type Pagination struct {
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}
