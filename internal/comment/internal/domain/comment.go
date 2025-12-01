package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	Id        bson.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Path      string
	Content   string
	RootId    bson.ObjectID
	ParentId  bson.ObjectID
	Nickname  string
	Avatar    string
	Email     string
	SiteUrl   string
	IsAdmin   bool
}

type CommentShow struct {
	Id        bson.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Path      string
	Content   string
	Nickname  string
	Avatar    string
	SiteUrl   string
	IsAdmin   bool
}
