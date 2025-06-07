package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Friend struct {
	ID bson.ObjectID
	Name string 
	Description string
	SiteUrl string
	WebsiteType string
	AvatarUrl string
	IsActive bool
}