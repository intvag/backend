package main

import (
	"github.com/gofrs/uuid/v5"
)

type Person struct {
	ID       uuid.UUID `json:"-" gorm:"primary_key"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Postcode string    `json:"postcode"`
}
