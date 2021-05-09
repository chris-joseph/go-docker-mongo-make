package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct{
	ID 	primitive.ObjectID 	`bson:"_id"`
	UserName string 		`json:"username"`
	Password string			`json:"-"`
}