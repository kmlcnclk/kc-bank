package domain

import (
	"time"
)

type Account struct {
	Id        string    `bson:"_id"`
	Currency  string    `bson:"currency" validate:"required"`
	Iban      string    `bson:"iban" validate:"required"`
	Balance   float64   `bson:"balance" validate:"required"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	UserId    string    `bson:"userId"`
}
