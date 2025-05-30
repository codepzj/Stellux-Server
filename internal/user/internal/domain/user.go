package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Nickname  string
	RoleId    int
	Avatar    string
	Email     string
}

type Page struct {
	PageNo   int64
	PageSize int64
}
