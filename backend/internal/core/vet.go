package core

import (
	"context"
	"time"
)

type Vets struct {
	ID                 int       `gorm:"column:id"`
	UserID             int       `gorm:"column:user_id"`
	IsOrganization     bool      `gorm:"column:is_organization"`
	Username           *string   `gorm:"column:username"`
	Firstname          *string   `gorm:"column:firstname"`
	Lastname           *string   `gorm:"column:lastname"`
	Patronymic         *string   `gorm:"column:patronymic"`
	Education          *string   `gorm:"column:education"`
	OrgName            *string   `gorm:"column:org_name"`
	Location           *string   `gorm:"column:location"`
	Price              float64   `gorm:"column:price"`
	OrgEmail           *string   `gorm:"column:org_email"`
	InnNumber          *string   `gorm:"column:inn_number"`
	IsRemoteConsulting bool      `gorm:"column:remote_consulting"`
	IsInpatient        bool      `gorm:"column:is_inpatient"`
	Description        *string   `gorm:"column:description"`
	IsDeleted          bool      `gorm:"column:is_deleted"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
	DeletedAt          time.Time `gorm:"column:deleted_at"`
}

type UpdateVets struct {
	ID                 int       `gorm:"column:id"`
	UserID             int       `gorm:"column:user_id"`
	Username           *string   `gorm:"column:username"`
	Firstname          *string   `gorm:"column:firstname"`
	Lastname           *string   `gorm:"column:lastname"`
	Patronymic         *string   `gorm:"column:patronymic"`
	Education          *string   `gorm:"column:education"`
	OrgName            *string   `gorm:"column:org_name"`
	Location           *string   `gorm:"column:location"`
	Price              float64   `gorm:"column:price"`
	OrgEmail           *string   `gorm:"column:org_email"`
	InnNumber          *string   `gorm:"column:Inn_number"`
	IsRemoteConsulting bool      `gorm:"column:remote_consulting"`
	IsOrganization     bool      `gorm:"column:is_organization"`
	IsInpatient        bool      `gorm:"column:is_inpatient"`
	Description        *string   `gorm:"column:description"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type VetsDetails struct {
	Vet  Vets
	User User
}

type GetAllVetParams struct {
	SortBy    *string
	SortOrder *string
	Location  *string
	MinRating *float64
	MaxRating *float64
	MinPrice  *float64
	MaxPrice  *float64
	Limit     *int
	Offset    *int
}

type VetStore interface {
	UpdateByID(ctx context.Context, update UpdateVets) (updatedVet Vets, err error)
	GetAll(ctx context.Context, params GetAllVetParams) ([]Vets, error)
	GetByUserID(ctx context.Context, userID int) (vet Vets, err error)
	GetByOrgName(ctx context.Context, orgName string) (vet Vets, err error)
	Create(ctx context.Context, vet Vets) error
	DeleteByUserID(ctx context.Context, userID int) error
}

type VetService interface {
	UpdateByUserID(ctx context.Context, vet UpdateVets) (VetsDetails, error)
	GetAll(ctx context.Context, params GetAllVetParams) ([]VetsDetails, error)
	GetByUserID(ctx context.Context, userID int) (vet VetsDetails, err error)
	GetByOrgName(ctx context.Context, orgName string) (vet VetsDetails, err error)
	Create(ctx context.Context, vet Vets) error
	DeleteByUserID(ctx context.Context, userID int) error

	VetReviewsService
}

// TableName table name in db for gorm
func (Vets) TableName() string {
	return "Vets"
}
