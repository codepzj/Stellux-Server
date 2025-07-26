package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Label struct {
	Id        bson.ObjectID
	LabelType string
	Name      string
}

type LabelPostCount struct {
	Id        bson.ObjectID
	LabelType string
	Name      string
	Count     int
}
