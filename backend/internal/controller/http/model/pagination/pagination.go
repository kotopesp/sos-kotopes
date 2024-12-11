package pagination

type Pagination struct {
	Total       int `json:"total" example:"10"`
	TotalPages  int `json:"total_pages" example:"10"`
	CurrentPage int `json:"current_page" example:"1"`
	PerPage     int `json:"per_page" example:"1"`
}
