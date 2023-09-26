package models

import (
	"backend/pkg/logger"
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseModel struct {
	db        *mongo.Database
	logger    *logger.DebugLogger
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type IBaseModel interface {
	Create()
	Delete()
	Find()
	GetAll()
	Update()
}

func Initialize(db *mongo.Database) *BaseModel {
	logmanager, err := logger.NewDebugLogger("database.log", os.Stdout, "BaseModel")
	if err != nil {
		panic(err)
	}
	return &BaseModel{
		db:     db,
		logger: logmanager,
	}

}

func (m *BaseModel) Create(data interface{}) {

	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	data_type := reflect.TypeOf(data)
	collection_name := data_type.Name()
	m.logger.Debug(fmt.Sprintf("Creating %v | %v", collection_name, data))
	err := m.db.CreateCollection(ctx, collection_name)
	if err != nil {
		log.Fatalf("Err:%v", err)
	}

	collection := m.db.Collection(collection_name)
	res, err := collection.InsertOne(ctx, data)

	m.ID = res.InsertedID.(primitive.ObjectID)

	if err != nil {
		log.Fatalf("Data didn't Insert into Collection.Err:%v", err)
	}

	defer cancel()
}

func (m *BaseModel) Find(data interface{}) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	data_type := reflect.TypeOf(data)
	collection_name := data_type.Name()
	m.logger.Debug(fmt.Sprintf("Finding %v | %v", collection_name, data))
	collection := m.db.Collection(collection_name)

	cursor, err := collection.Find(ctx, data)

	if err != nil {
		log.Fatalf("Cursor Error.Err:%v", err)
	}

	res := new(interface{})
	err = cursor.All(ctx, res)
	if err != nil {
		panic(err)
	}
	defer cancel()

	return res
}
