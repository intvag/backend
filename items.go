package main

import (
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model

	ID              uuid.UUID
	Category        string
	Manufacturer    string
	ModelName       string
	AverageLifetime int
	Lastability     float64
	Repairability   float64
}
