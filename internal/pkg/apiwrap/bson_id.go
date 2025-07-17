package apiwrap

import (
	"encoding/json"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BsonId bson.ObjectID

// 自定义 JSON 解析逻辑
func (b *BsonId) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return errors.Wrapf(err, "json unmarshal failed: %s", data)
	}

	objId, err := bson.ObjectIDFromHex(str)
	if err != nil {
		return errors.Wrapf(err, "bsonid from hex failed: %s", str)
	}
	*b = BsonId(objId)
	return nil
}

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
