package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Label struct {
	ID        bson.ObjectID `bson:"_id"`
	LabelType string        `bson:"label_type"`
	Name      string        `bson:"name"`
}
