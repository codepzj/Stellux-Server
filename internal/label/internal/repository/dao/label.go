package dao

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Label struct {
	ID        bson.ObjectID `bson:"_id"`
	LabelType string        `bson:"type"`
	Name      string        `bson:"name"`
}

// LabelPostCount 每种标签文章的数量
type LabelPostCount struct {
	ID        bson.ObjectID `bson:"_id"`
	LabelType string        `bson:"type"`
	Name      string        `bson:"name"`
	Count     int           `bson:"post_count"`
}

type ILabelDao interface {
	CreateLabel(ctx context.Context, label *Label) error
	UpdateLabel(ctx context.Context, id bson.ObjectID, label *Label) error
	DeleteLabel(ctx context.Context, id bson.ObjectID) error
	GetLabelById(ctx context.Context, id bson.ObjectID) (*Label, error)
	QueryLabelList(ctx context.Context, labelType string, keyword string, limit int64, skip int64) ([]*Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error)
	GetCategoryLabelWithCount(ctx context.Context) ([]*LabelPostCount, error)
	GetTagsLabelWithCount(ctx context.Context) ([]*LabelPostCount, error)
	GetLabelByName(ctx context.Context, name string) (*Label, error)
}

var _ ILabelDao = (*LabelDao)(nil)

func NewLabelDao(db *mongo.Database) *LabelDao {
	return &LabelDao{coll: db.Collection("label")}
}

type LabelDao struct {
	coll *mongo.Collection
}

// CreateLabel 创建标签
func (d *LabelDao) CreateLabel(ctx context.Context, label *Label) error {
	label.ID = bson.NewObjectID()
	_, err := d.coll.InsertOne(ctx, label)
	return err
}

// UpdateLabel 更新标签
func (d *LabelDao) UpdateLabel(ctx context.Context, id bson.ObjectID, label *Label) error {
	update := bson.M{
		"$set": bson.M{
			"name": label.Name,
			"type": label.LabelType,
		},
	}
	_, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// DeleteLabel 删除标签
func (d *LabelDao) DeleteLabel(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// GetLabelById 根据id获取标签
func (d *LabelDao) GetLabelById(ctx context.Context, id bson.ObjectID) (*Label, error) {
	var label Label
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}

// QueryLabelList 分页查询标签
func (d *LabelDao) QueryLabelList(ctx context.Context, labelType string, keyword string, limit int64, skip int64) ([]*Label, int64, error) {
	filter := bson.M{"type": labelType}
	if keyword != "" {
		filter["name"] = bson.M{"$regex": keyword, "$options": "i"}
	}

	count, err := d.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := d.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var labelList []*Label
	if err = cursor.All(ctx, &labelList); err != nil {
		return nil, 0, err
	}
	return labelList, count, nil
}

// GetAllLabelsByType 通过类型获取所有标签
func (d *LabelDao) GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error) {
	cursor, err := d.coll.Find(ctx, bson.M{"type": labelType})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var labels []*Label
	if err = cursor.All(ctx, &labels); err != nil {
		return nil, err
	}
	return labels, nil
}

// GetCategoryLabelWithCount 获取分类标签及其文章数量, 一篇文章只能有一个分类
func (d *LabelDao) GetCategoryLabelWithCount(ctx context.Context) ([]*LabelPostCount, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "post"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "category_id"},
			{Key: "as", Value: "category_post"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "post_count", Value: bson.D{{Key: "$size", Value: "$category_post"}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "post_count", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "post_count", Value: -1}}}},
	}

	cursor, err := d.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var labelPostCount []*LabelPostCount
	if err = cursor.All(ctx, &labelPostCount); err != nil {
		return nil, err
	}
	return labelPostCount, nil
}

// GetTagsLabelWithCountByType 获取标签及其文章数量, 一篇文章可以有多个标签
func (d *LabelDao) GetTagsLabelWithCount(ctx context.Context) ([]*LabelPostCount, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "post"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "tags_id"},
			{Key: "as", Value: "tags_post"},
		}}},
		{{Key: "$addFields", Value: bson.D{
			{Key: "post_count", Value: bson.D{{Key: "$size", Value: "$tags_post"}}},
		}}},
		{{Key: "$match", Value: bson.D{
			{Key: "post_count", Value: bson.D{{Key: "$gt", Value: 0}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "post_count", Value: -1}}}},
	}

	cursor, err := d.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var labelPostCount []*LabelPostCount
	if err = cursor.All(ctx, &labelPostCount); err != nil {
		return nil, err
	}
	return labelPostCount, nil
}

// GetLabelByName 根据名称获取标签
func (d *LabelDao) GetLabelByName(ctx context.Context, name string) (*Label, error) {
	var label Label
	err := d.coll.FindOne(ctx, bson.M{"name": name}).Decode(&label)
	if err != nil {
		return nil, err
	}
	return &label, nil
}
