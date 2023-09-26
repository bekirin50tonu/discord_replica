package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/internal/user/models"
)

// BookRepository MongoDB işlemleri için bir yapı
type AccountRepository struct {
	collection *mongo.Collection
}

// NewBookRepository yeni bir BookRepository oluşturur
func NewAccountRepository(db *mongo.Database, collectionName string) (*AccountRepository, error) {
	err := db.CreateCollection(context.Background(), collectionName)
	if err != nil {
		return nil, err
	}
	collection := db.Collection(collectionName) // Koleksiyon adınıza uygun olarak değiştirin
	return &AccountRepository{
		collection: collection,
	}, nil
}

// Create yeni bir kitap ekler
func (r *AccountRepository) Create(account *models.Account) (*models.Account, error) {
	res, err := r.collection.InsertOne(context.TODO(), account)
	if err != nil {
		return nil, err
	}
	account.ID = res.InsertedID.(primitive.ObjectID)

	return account, nil
}

// GetAll tüm kitapları getirir
func (r *AccountRepository) GetAll() ([]*models.Account, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var books []*models.Account
	for cur.Next(context.Background()) {
		var book models.Account
		err := cur.Decode(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	return books, nil
}

func (r *AccountRepository) Delete(accountID primitive.ObjectID) error {
	filter := bson.M{"_id": accountID}
	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Account Not Found.")
	}
	return nil
}

// FindOne belirli bir hesabı ID ile bulur
func (r *AccountRepository) FindOne(accountID primitive.ObjectID) (*models.Account, error) {
	filter := bson.M{"_id": accountID}
	var account models.Account
	err := r.collection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// FindAll tüm hesapları getirir
func (r *AccountRepository) FindAll() ([]*models.Account, error) {
	cur, err := r.collection.Find(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	var accounts []*models.Account
	for cur.Next(context.Background()) {
		var account models.Account
		err := cur.Decode(&account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

// Update bir hesabın bilgilerini günceller
func (r *AccountRepository) Update(accountID primitive.ObjectID, updatedAccount *models.Account) error {
	filter := bson.M{"_id": accountID}
	update := bson.M{
		"$set": updatedAccount,
	}
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("Account Not Found.")
	}
	return nil
}
