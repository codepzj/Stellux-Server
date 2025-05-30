package apiwrap

import (
	"log/slog"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BsonID string

func (bid BsonID) IsZero() bool {
	return bid == "" || bid == "000000000000000000000000"
}

func ConvertBsonID(id string) bson.ObjectID {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		slog.Warn("ConvertBsonID", "error", err.Error())
		return bson.ObjectID{}
	}
	return objID
}

func ConvertBsonIDList(idList []string) []bson.ObjectID {
	return lo.Map(idList, func(id string, _ int) bson.ObjectID {
		return ConvertBsonID(id)
	})
}
