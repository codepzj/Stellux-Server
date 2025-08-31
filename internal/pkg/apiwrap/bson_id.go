package apiwrap

import (
	"log/slog"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BsonId bson.ObjectID

func (b BsonId) ToObjectID() bson.ObjectID {
	return bson.ObjectID(b)
}

func ConvertBsonID(id string) BsonId {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		slog.Warn("ConvertBsonID", "error", err.Error())
		return BsonId{}
	}
	return BsonId(objId)
}

func ConvertBsonIDList(idList []string) []BsonId {
	return lo.Map(idList, func(id string, _ int) BsonId {
		return ConvertBsonID(id)
	})
}

func ToObjectIDList(idList []BsonId) []bson.ObjectID {
	return lo.Map(idList, func(id BsonId, _ int) bson.ObjectID {
		return id.ToObjectID()
	})
}
