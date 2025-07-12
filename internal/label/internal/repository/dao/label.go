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
}

var _ ILabelDao = (*LabelDao)(nil)

func NewLabelDao(db *mongox.Database) *LabelDao {
	return &LabelDao{coll: mongox.NewCollection[Label](db, "label")}
}

type LabelDao struct {
	coll *mongox.Collection[Label]
}

func (d *LabelDao) CreateLabel(ctx context.Context, label *Label) error {
	label.ID = bson.NewObjectID()
	_, err := d.coll.Creator().InsertOne(ctx, label)
	return err
}

func (d *LabelDao) UpdateLabel(ctx context.Context, id bson.ObjectID, label *Label) error {
	_, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.SetFields(map[string]any{
		"name": label.Name,
		"type": label.LabelType,
	})).UpdateOne(ctx)
	return err
}

func (d *LabelDao) DeleteLabel(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	return err
}

func (d *LabelDao) GetLabelById(ctx context.Context, id bson.ObjectID) (*Label, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

// QueryLabelList 分页查询
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

// GetCategoryLabelWithCount 获取分类标签及其文章数量
func (d *LabelDao) GetCategoryLabelWithCount(ctx context.Context) ([]*LabelPostCount, error) {

	// 标签根据文章数量降序排序
	aggregationBuilder := aggregation.NewStageBuilder().Lookup("post", "category_post", &aggregation.LookUpOptions{
		LocalField:   "_id",
		ForeignField: "category_id",
	}).AddFields(bson.M{
		"post_count": bson.M{
			"$size": "$category_post",
		}})

	var labelPostCount []*LabelPostCount

	err := d.coll.Aggregator().Pipeline(aggregationBuilder.Build()).AggregateWithParse(ctx, &labelPostCount)
	if err != nil {
		return nil, err
	}
	return labelPostCount, nil
}
