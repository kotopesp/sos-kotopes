package entity

type (
	Entity struct {
		Field1 int    `json:"field1"`
		Field2 string `json:"field2"`
	}

	GetAllParams struct {
		Limit      int    `query:"limit"`
		Offset     int    `query:"offset"`
		Sort       string `query:"sort"`
		SearchTerm string `query:"q"`
		Status     string `query:"status"`
	}
)
