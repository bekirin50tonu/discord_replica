package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Server struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"server_name"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
