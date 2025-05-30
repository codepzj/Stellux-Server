package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type File struct {
	ID       bson.ObjectID
	FileName string
	Url      string
	Dst      string
}
