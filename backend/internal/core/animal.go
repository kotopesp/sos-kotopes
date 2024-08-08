package core

import (
	"context"
	"time"
)

type (
	Animal struct {
		ID          int       `gorm:"column:id;primaryKey"` // Unique identifier for the animal
		KeeperID    int       `gorm:"column:keeper_id"`     // ID of the person responsible for the animal
		AnimalType  string    `gorm:"column:animal_type"`   // Type of the animal (dog, cat)
		Age         int       `gorm:"column:age"`           // Age of the animal
		Color       string    `gorm:"column:color"`         // Color of the animal
		Gender      string    `gorm:"column:gender"`        // Gender of the animal (male, female)
		Description string    `gorm:"column:description"`   // Description of the animal
		Status      string    `gorm:"column:status"`        // Status of the animal (lost, found, need home)
		CreatedAt   time.Time `gorm:"column:created_at"`    // Timestamp when the record was created
		UpdatedAt   time.Time `gorm:"column:updated_at"`    // Timestamp when the record was last updated
	}

	UpdateRequestBodyAnimal struct {
		AnimalType  string `gorm:"column:animal_type"`
		Age         int    `gorm:"column:age"`
		Color       string `gorm:"column:color"`
		Gender      string `gorm:"column:gender"`
		Description string `gorm:"column:description"`
		Status      string `gorm:"column:status"`
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
