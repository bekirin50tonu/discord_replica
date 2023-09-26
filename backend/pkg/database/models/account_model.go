package models

type Account struct {
	BaseModel
	Username string `bson:"username,omitempty"`
	Password string `bson:"password,omitempty"`
}

func NewAccount(username string, password string) *Account {
	return &Account{Username: username, Password: password}
}
