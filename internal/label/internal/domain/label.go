package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Label struct {
	ID        bson.ObjectID
	LabelType string
	Name      string
}

type LabelPostCount struct {
	ID        bson.ObjectID
	LabelType string
	Name      string
	Count     int
}
