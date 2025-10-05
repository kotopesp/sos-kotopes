package report

type (
	// CreateRequestBodyReport - represent body of create report request.
	CreateRequestBodyReport struct {
		TargetID   int    `json:"target_id" validate:"required,min=1"`
		TargetType string `json:"target_type" validate:"required,oneof=post comment"`
		Reason     string `json:"reason" validate:"required,oneof=spam violent_content violent_speech"`
	}
)
