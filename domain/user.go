package domain

import (
	"time"
)

type User struct {
	Id        string    `bson:"_id"`
	FirstName string    `bson:"firstName" validate:"required"`
	LastName  string    `bson:"lastName" validate:"required"`
	Email     string    `bson:"email" validate:"required,email"`
	Password  string    `bson:"password" validate:"required,min=6"`
	Age       int32     `bson:"age" validate:"gte=0,lte=130"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
