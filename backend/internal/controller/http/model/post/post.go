package post

type (
	Post struct {
		ID        int    `gorm:"primary key;autoIncrement" json:"id"`
		Title     string `json:"title"`
		Body      string `json:"body"`
		UserID    int    `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		AnimalID  int    `json:"animal_id"`
	}
	GetAllPostsParams struct {
		Limit      int    `query:"limit"`
		Offset     int    `query:"offset"`
		SortBy     string `query:"sort_by"`
		SortOrder  string `query:"sort_order"`
		SearchTerm string `query:"q"`
	}
)
