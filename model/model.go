package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty" validate:"required,min=2,max=40"`
	Age          int                `json:"age,omitempty" bson:"age,omitempty" validate:"required"`
	Code         int                `json:"code,omitempty" bson:"code,omitempty"`
	CodeID       string             `json:"code_id,omitempty" bson:"code_id,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty" validate:"email,required"`
	Password     string             `json:"password,omitempty" bson:"password,omitempty" validate:"required"`
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	IpAddress    string             `json:"ip_address,omitempty" bson:"ip_address,omitempty"`
	RegisteredAt time.Time          `json:"registered_at,omitempty" bson:"registered_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
