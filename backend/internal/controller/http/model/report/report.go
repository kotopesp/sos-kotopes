package report

type (
	// CreateRequestBodyReport - represent body of create report request.
	CreateRequestBodyReport struct {
		Reason string `json:"reason" validate:"required,oneof=spam violent_speech violent_content"`
	}
)
