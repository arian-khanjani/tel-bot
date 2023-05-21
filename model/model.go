package model

type User struct {
	ID        int64   `json:"id" bson:"_id"`
	Username  string  `json:"username" bson:"username"`
	FirstName string  `json:"first_name" bson:"first_name"`
	Account   Account `json:"account" bson:"account"`
}

type Account struct {
	Balance int `json:"balance" bson:"balance"`
}

type Client struct {
	ID         int64  `json:"id" bson:"_id"`
	ProviderID int64  `json:"provider_id" bson:"provider_id"`
	Username   string `json:"username" bson:"username"`
}
