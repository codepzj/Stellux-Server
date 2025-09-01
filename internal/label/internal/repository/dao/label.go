package dao

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	QueryLabelList(ctx context.Context, labelType string, limit int64, skip int64) ([]*Label, int64, error)
	GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error)
	GetCategoryLabelWithCount(ctx context.Context) ([]*LabelPostCount, error)
	GetTagsLabelWithCount(ctx context.Context) ([]*LabelPostCount, error)
	GetLabelByName(ctx context.Context, name string) (*Label, error)
}

var _ ILabelDao = (*LabelDao)(nil)

func NewLabelDao(db *mongox.Database) *LabelDao {
	return &LabelDao{coll: mongox.NewCollection[Label](db, "label")}
}

type LabelDao struct {
	coll *mongox.Collection[Label]
}

// CreateLabel 创建标签
func (d *LabelDao) CreateLabel(ctx context.Context, label *Label) error {
	label.ID = bson.NewObjectID()
	_, err := d.coll.Creator().InsertOne(ctx, label)
	return err
}

// UpdateLabel 更新标签
func (d *LabelDao) UpdateLabel(ctx context.Context, id bson.ObjectID, label *Label) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(map[string]any{
		"name": label.Name,
		"type": label.LabelType,
	})).UpdateOne(ctx)
	return err
}

// DeleteLabel 删除标签
func (d *LabelDao) DeleteLabel(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	return err
}

// GetLabelById 根据id获取标签
func (d *LabelDao) GetLabelById(ctx context.Context, id bson.ObjectID) (*Label, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

// QueryLabelList 分页查询标签
func (d *LabelDao) QueryLabelList(ctx context.Context, labelType string, limit int64, skip int64) ([]*Label, int64, error) {
	count, err := d.coll.Finder().Filter(query.Eq("type", labelType)).Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	labelList, err := d.coll.Finder().Filter(query.Eq("type", labelType)).Limit(limit).Skip(skip).Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return labelList, count, nil
}

// GetAllLabelsByType 通过类型获取所有标签
func (d *LabelDao) GetAllLabelsByType(ctx context.Context, labelType string) ([]*Label, error) {
	return d.coll.Finder().Filter(query.Eq("type", labelType)).Find(ctx)
}

// GetCategoryLabelWithCount 获取分类标签及其文章数量, 一篇文章只能有一个分类
func (d *LabelDao) GetCategoryLabelWithCount(ctx context.Context) ([]*LabelPostCount, error) {

	// 标签根据文章数量降序排序, 并且筛选出文章数量大于0的标签
	aggregationBuilder := aggregation.NewStageBuilder().
		Lookup("post", "category_post", &aggregation.LookUpOptions{
			LocalField:   "_id",
			ForeignField: "category_id",
		}).
		AddFields(bson.M{
			"post_count": bson.M{
				"$size": "$category_post",
			}}).
		Match(bson.M{
			"post_count": bson.M{
				"$gt": 0,
			},
		}).
		Sort(bson.M{
			"post_count": -1,
		})

	var labelPostCount []*LabelPostCount

	err := d.coll.Aggregator().Pipeline(aggregationBuilder.Build()).AggregateWithParse(ctx, &labelPostCount)
	if err != nil {
		return nil, err
	}
	return labelPostCount, nil
}

// GetTagsLabelWithCountByType 获取标签及其文章数量, 一篇文章可以有多个标签
func (d *LabelDao) GetTagsLabelWithCount(ctx context.Context) ([]*LabelPostCount, error) {
	aggregationBuilder := aggregation.NewStageBuilder().
		Lookup("post", "tags_post", &aggregation.LookUpOptions{
			LocalField:   "_id",
			ForeignField: "tags_id",
		}).
		AddFields(bson.M{
			"post_count": bson.M{
				"$size": "$tags_post",
			}}).
		Match(bson.M{
			"post_count": bson.M{
				"$gt": 0,
			},
		}).
		Sort(bson.M{
			"post_count": -1,
		})

	var labelPostCount []*LabelPostCount

	err := d.coll.Aggregator().Pipeline(aggregationBuilder.Build()).AggregateWithParse(ctx, &labelPostCount)
	if err != nil {
		return nil, err
	}
	return labelPostCount, nil
}

// GetLabelByName 根据名称获取标签
func (d *LabelDao) GetLabelByName(ctx context.Context, name string) (*Label, error) {
	return d.coll.Finder().Filter(query.Eq("name", name)).FindOne(ctx)
}