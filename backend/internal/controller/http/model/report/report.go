package report

type (
	CreateRequestBodyReport struct {
		UserID int    `json:"user_id"`
		PostID int    `json:"post_id"`
		Reason string `json:"reason"`
	}
)
