package core

import (
	"time"
	"context"
)

type (
	Animal struct {
		ID          int       `gorm:"column:id;primaryKey"`
		KeeperID    int       `gorm:"column:keeper_id"`
		AnimalType  string    `gorm:"column:animal_type"`
		Age         int       `gorm:"column:age"`
		Color       string    `gorm:"column:color"`
		Gender      string    `gorm:"column:gender"`
		Description string    `gorm:"column:description"`
		Status      string    `gorm:"column:status"`
		CreatedAt   time.Time `gorm:"column:created_at"`
		UpdatedAt   time.Time `gorm:"column:updated_at"`
	}

	UpdateRequestBodyAnimal struct {
		AnimalType  string    `gorm:"column:animal_type"`
		Age         int       `gorm:"column:age"`
		Color       string    `gorm:"column:color"`
		Gender      string    `gorm:"column:gender"`
		Description string    `gorm:"column:description"`
		Status      string    `gorm:"column:status"`
	}

	AnimalStore interface {
		CreateAnimal(ctx context.Context, animal Animal) (Animal, error)
		GetAnimalByID(ctx context.Context, id int) (Animal, error)
		UpdateAnimal(ctx context.Context, animal Animal) (Animal, error)
	}
)

func (Animal) TableName() string {
	return "animals"
}
