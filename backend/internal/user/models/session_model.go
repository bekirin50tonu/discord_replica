package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id,omitempty"`

	AccessToken  string `bson:"access_token"`
	RefreshToken string `bson:"refresh_token"`

	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}
