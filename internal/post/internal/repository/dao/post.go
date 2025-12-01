package dao

import (
	"context"
	"time"

	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/pkg/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Post struct {
	ID          bson.ObjectID   `bson:"_id,omitempty"`
	CreatedAt   time.Time       `bson:"created_at"`
	UpdatedAt   time.Time       `bson:"updated_at"`
	DeletedAt   *time.Time      `bson:"deleted_at,omitempty"`
	Title       string          `bson:"title"`
	Content     string          `bson:"content"`
	Description string          `bson:"description"`
	Author      string          `bson:"author"`
	Alias       string          `bson:"alias"`
	CategoryID  bson.ObjectID   `bson:"category_id"`
	TagsID      []bson.ObjectID `bson:"tags_id"`
	IsPublish   bool            `bson:"is_publish"`
	IsTop       bool            `bson:"is_top"`
	Thumbnail   string          `bson:"thumbnail"`
}

type PostUpdate struct {
	Title       string          `bson:"title"`
	Content     string          `bson:"content"`
	Description string          `bson:"description"`
	Author      string          `bson:"author"`
	Alias       string          `bson:"alias"`
	CategoryID  bson.ObjectID   `bson:"category_id"`
	TagsID      []bson.ObjectID `bson:"tags_id"`
	IsPublish   bool            `bson:"is_publish"`
	IsTop       bool            `bson:"is_top"`
	Thumbnail   string          `bson:"thumbnail"`
}

// 聚合查询返回带有category和tags的结构体
type PostCategoryTags struct {
	Id          bson.ObjectID  `bson:"_id"`
	CreatedAt   time.Time      `bson:"created_at"`
	UpdatedAt   time.Time      `bson:"updated_at"`
	Title       string         `bson:"title"`
	Content     string         `bson:"content"`
	Description string         `bson:"description"`
	Author      string         `bson:"author"`
	Alias       string         `bson:"alias"`
	Category    label.Domain   `bson:"category"`
	Tags        []label.Domain `bson:"tags"`
	IsPublish   bool           `bson:"is_publish"`
	IsTop       bool           `bson:"is_top"`
	Thumbnail   string         `bson:"thumbnail"`
}

type UpdatePost struct {
	CreatedAt   time.Time       `bson:"created_at,omitempty"`
	Title       string          `bson:"title"`
	Content     string          `bson:"content"`
	Description string          `bson:"description"`
	Author      string          `bson:"author"`
	Alias       string          `bson:"alias"`
	CategoryId  bson.ObjectID   `bson:"category_id"`
	TagsId      []bson.ObjectID `bson:"tags_id"`
	IsPublish   bool            `bson:"is_publish"`
	IsTop       bool            `bson:"is_top"`
	Thumbnail   string          `bson:"thumbnail"`
}

type IPostDao interface {
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, id bson.ObjectID, post *UpdatePost) error
	UpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	SoftDelete(ctx context.Context, id bson.ObjectID) error
	SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Restore(ctx context.Context, id bson.ObjectID) error
	RestoreBatch(ctx context.Context, ids []bson.ObjectID) error
	GetByID(ctx context.Context, id bson.ObjectID) (*Post, error)
	GetByKeyWord(ctx context.Context, keyWord string) ([]*Post, error)
	GetDetailByID(ctx context.Context, id bson.ObjectID) (*PostCategoryTags, error)
	GetList(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D) ([]*PostCategoryTags, int64, error)
	GetListWithTagFilter(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D, hasTagFilter bool, labelName string) ([]*PostCategoryTags, int64, error)
	GetListWithFilter(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D, hasTagFilter bool, labelName string, hasCategoryFilter bool, categoryName string) ([]*PostCategoryTags, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*Post, error)
	FindByAlias(ctx context.Context, alias string) (*Post, error)
}

var _ IPostDao = (*PostDao)(nil)

func NewPostDao(db *mongo.Database) *PostDao {
	return &PostDao{coll: db.Collection("post")}
}

type PostDao struct {
	coll *mongo.Collection
}

// Create 创建文章
func (d *PostDao) Create(ctx context.Context, post *Post) error {
	post.ID = bson.NewObjectID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	insertResult, err := d.coll.InsertOne(ctx, post)
	if err != nil {
		return err
	}
	if insertResult.InsertedID == nil {
		return errors.Wrap(err, "插入文章失败")
	}
	return nil
}

func (d *PostDao) Update(ctx context.Context, id bson.ObjectID, post *UpdatePost) error {
	update := bson.M{"$set": post}
	if update["$set"].(bson.M)["updated_at"] == nil {
		update["$set"].(bson.M)["updated_at"] = time.Now()
	}
	updateResult, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if updateResult.ModifiedCount == 0 {
		return errors.New("文章修改失败")
	}
	return nil
}

// UpdatePostPublishStatus 更新文章发布状态
func (d *PostDao) UpdatePostPublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	update := bson.M{"$set": bson.M{"is_publish": isPublish, "updated_at": time.Now()}}
	_, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// SoftDelete 软删除文章
func (d *PostDao) SoftDelete(ctx context.Context, id bson.ObjectID) error {
	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now, "is_publish": false, "is_top": false, "updated_at": now}}
	_, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete 删除文章
func (d *PostDao) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := d.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteBatch 批量删除文章
func (d *PostDao) DeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	_, err := d.coll.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	return err
}

// Restore 恢复文章
func (d *PostDao) Restore(ctx context.Context, id bson.ObjectID) error {
	update := bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := d.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetDetailByID 获取文章
func (d *PostDao) GetDetailByID(ctx context.Context, id bson.ObjectID) (*PostCategoryTags, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "label"},
			{Key: "localField", Value: "category_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: "$category"}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "label"},
			{Key: "localField", Value: "tags_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tags"},
		}}},
	}

	cursor, err := d.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var postResult []PostCategoryTags
	if err = cursor.All(ctx, &postResult); err != nil {
		return nil, err
	}
	if len(postResult) == 0 {
		return nil, errors.New("文章不存在")
	}
	return &postResult[0], nil
}

// GetByID 获取文章
func (d *PostDao) GetByID(ctx context.Context, id bson.ObjectID) (*Post, error) {
	var post Post
	err := d.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetByKeyWord 获取文章
func (d *PostDao) GetByKeyWord(ctx context.Context, keyWord string) ([]*Post, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": keyWord, "$options": "i"}},
			{"description": bson.M{"$regex": keyWord, "$options": "i"}},
		},
		"deleted_at": nil,
		"is_publish": true,
	}
	cursor, err := d.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetList 获取文章列表
func (d *PostDao) GetList(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D) ([]*PostCategoryTags, int64, error) {
	cursor, err := d.coll.Aggregate(ctx, pagePipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var postResult []PostCategoryTags
	if err = cursor.All(ctx, &postResult); err != nil {
		return nil, 0, err
	}

	count, err := d.coll.CountDocuments(ctx, cond)
	if err != nil {
		return nil, 0, err
	}

	return utils.ValToPtrList(postResult), count, nil
}

// buildCountPipeline 构建计数管道
func (d *PostDao) buildCountPipeline(cond bson.D, labelName string, categoryName string) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: cond}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "label"},
			{Key: "localField", Value: "category_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "category"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$category"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "label"},
			{Key: "localField", Value: "tags_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "tags"},
		}}},
	}

	if labelName != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "type", Value: "tag"},
				{Key: "name", Value: labelName},
			}}}},
		}}})
	}

	if categoryName != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "category.type", Value: "category"},
			{Key: "category.name", Value: categoryName},
		}}})
	}

	return pipeline
}

// GetListWithTagFilter 获取文章列表（带标签过滤）
func (d *PostDao) GetListWithTagFilter(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D, hasTagFilter bool, labelName string) ([]*PostCategoryTags, int64, error) {
	cursor, err := d.coll.Aggregate(ctx, pagePipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var postResult []PostCategoryTags
	if err = cursor.All(ctx, &postResult); err != nil {
		return nil, 0, err
	}

	var count int64
	if hasTagFilter {
		countPipeline := d.buildCountPipeline(cond, labelName, "")
		countCursor, err := d.coll.Aggregate(ctx, countPipeline)
		if err != nil {
			return nil, 0, err
		}
		defer countCursor.Close(ctx)

		var countResult []bson.M
		if err = countCursor.All(ctx, &countResult); err != nil {
			return nil, 0, err
		}
		count = int64(len(countResult))
	} else {
		count, err = d.coll.CountDocuments(ctx, cond)
		if err != nil {
			return nil, 0, err
		}
	}

	return utils.ValToPtrList(postResult), count, nil
}

// GetListWithFilter 获取文章列表（带标签和分类过滤）
func (d *PostDao) GetListWithFilter(ctx context.Context, pagePipeline mongo.Pipeline, cond bson.D, hasTagFilter bool, labelName string, hasCategoryFilter bool, categoryName string) ([]*PostCategoryTags, int64, error) {
	cursor, err := d.coll.Aggregate(ctx, pagePipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var postResult []PostCategoryTags
	if err = cursor.All(ctx, &postResult); err != nil {
		return nil, 0, err
	}

	var count int64
	if hasTagFilter || hasCategoryFilter {
		countPipeline := d.buildCountPipeline(cond, labelName, categoryName)
		countCursor, err := d.coll.Aggregate(ctx, countPipeline)
		if err != nil {
			return nil, 0, err
		}
		defer countCursor.Close(ctx)

		var countResult []bson.M
		if err = countCursor.All(ctx, &countResult); err != nil {
			return nil, 0, err
		}
		count = int64(len(countResult))
	} else {
		count, err = d.coll.CountDocuments(ctx, cond)
		if err != nil {
			return nil, 0, err
		}
	}

	return utils.ValToPtrList(postResult), count, nil
}

// SoftDeleteBatch 批量软删除文章
func (d *PostDao) SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	now := time.Now()
	update := bson.M{"$set": bson.M{"deleted_at": now, "is_publish": false, "is_top": false, "updated_at": now}}
	_, err := d.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, update)
	return err
}

// RestoreBatch 批量恢复文章
func (d *PostDao) RestoreBatch(ctx context.Context, ids []bson.ObjectID) error {
	update := bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": time.Now()}}
	_, err := d.coll.UpdateMany(ctx, bson.M{"_id": bson.M{"$in": ids}}, update)
	return err
}

// GetAllPublishPost 获取所有发布文章
func (d *PostDao) GetAllPublishPost(ctx context.Context) ([]*Post, error) {
	opts := options.Find().SetSort(bson.M{"updated_at": -1})
	cursor, err := d.coll.Find(ctx, bson.M{"is_publish": true}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByAlias 根据别名获取文章
func (d *PostDao) FindByAlias(ctx context.Context, alias string) (*Post, error) {
	var post Post
	err := d.coll.FindOne(ctx, bson.M{"alias": alias}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
