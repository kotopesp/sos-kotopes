package core

import "context"

type (
	Entity struct {
		Field1 int    `db:"field1"`
		Field2 string `db:"field2"`
	}

	EntityStore interface {
		GetAll(ctx context.Context, params GetAllParams) (data []Entity, err error)
		GetByID(ctx context.Context, id int) (data Entity, err error)
	}

	EntityService interface {
		GetAll(ctx context.Context, params GetAllParams) (data []Entity, total int, err error)
		GetByID(ctx context.Context, id int) (data Entity, err error)
	}

	GetAllParams struct {
		SortBy     *string
		SortOrder  *string
		SearchTerm *string
		Limit      *int
		Offset     *int
	}
)
