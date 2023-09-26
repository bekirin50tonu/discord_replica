package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Channel struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"room_name"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
