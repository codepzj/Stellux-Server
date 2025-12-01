package repository

import (
	"context"

	"github.com/codepzj/Stellux-Server/internal/post/internal/domain"
	"github.com/codepzj/Stellux-Server/internal/post/internal/repository/dao"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IPostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	Update(ctx context.Context, post *domain.Post) error
	UpdatePublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error
	SoftDelete(ctx context.Context, id bson.ObjectID) error
	SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Delete(ctx context.Context, id bson.ObjectID) error
	DeleteBatch(ctx context.Context, ids []bson.ObjectID) error
	Restore(ctx context.Context, id bson.ObjectID) error
	RestoreBatch(ctx context.Context, ids []bson.ObjectID) error
	GetByID(ctx context.Context, id bson.ObjectID) (*domain.Post, error)
	GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error)
	GetDetailByID(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error)
	GetList(ctx context.Context, page *domain.PostQueryPage, postType string) ([]*domain.PostDetail, int64, error)
	GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error)
	FindByAlias(ctx context.Context, alias string) (*domain.Post, error)
}

var _ IPostRepository = (*PostRepository)(nil)

func NewPostRepository(dao dao.IPostDao) *PostRepository {
	return &PostRepository{dao: dao}
}

type PostRepository struct {
	dao dao.IPostDao
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.dao.Create(ctx, r.PostDomainToPostDO(post))
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.dao.Update(ctx, post.Id, r.PostDomainToUpdatePostDO(post))
}

// UpdatePostPublishStatus 更新文章发布状态
func (r *PostRepository) UpdatePublishStatus(ctx context.Context, id bson.ObjectID, isPublish bool) error {
	return r.dao.UpdatePostPublishStatus(ctx, id, isPublish)
}

func (r *PostRepository) SoftDelete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.SoftDelete(ctx, id)
}

// SoftDeleteBatch 批量软删除文章
func (r *PostRepository) SoftDeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.SoftDeleteBatch(ctx, ids)
}

// Delete 删除文章
func (r *PostRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Delete(ctx, id)
}

// DeleteBatch 批量删除文章
func (r *PostRepository) DeleteBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.DeleteBatch(ctx, ids)
}

// Restore 恢复文章
func (r *PostRepository) Restore(ctx context.Context, id bson.ObjectID) error {
	return r.dao.Restore(ctx, id)
}

// RestoreBatch 批量恢复文章
func (r *PostRepository) RestoreBatch(ctx context.Context, ids []bson.ObjectID) error {
	return r.dao.RestoreBatch(ctx, ids)
}

// GetByID 获取文章
func (r *PostRepository) GetByID(ctx context.Context, id bson.ObjectID) (*domain.Post, error) {
	post, err := r.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomain(post), nil
}

// GetDetailByID 获取文章
func (r *PostRepository) GetDetailByID(ctx context.Context, id bson.ObjectID) (*domain.PostDetail, error) {
	postCategoryTags, err := r.dao.GetDetailByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.PostCategoryTagsDOToPostDetail(postCategoryTags), nil
}

// GetByKeyWord 获取文章
func (r *PostRepository) GetByKeyWord(ctx context.Context, keyWord string) ([]*domain.Post, error) {
	posts, err := r.dao.GetByKeyWord(ctx, keyWord)
	if err != nil {
		return nil, err
	}
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	}), nil
}

// buildPostQueryCondition 构建文章查询条件
func (r *PostRepository) buildPostQueryCondition(page *domain.PostQueryPage, postType string) bson.D {
	conditions := bson.D{}

	if page.Keyword != "" {
		conditions = append(conditions, bson.E{Key: "$or", Value: []bson.M{
			{"title": bson.M{"$regex": page.Keyword, "$options": "i"}},
			{"description": bson.M{"$regex": page.Keyword, "$options": "i"}},
		}})
	}

	switch postType {
	case "publish":
		conditions = append(conditions, bson.E{Key: "deleted_at", Value: nil})
		conditions = append(conditions, bson.E{Key: "is_publish", Value: true})
	case "draft":
		conditions = append(conditions, bson.E{Key: "deleted_at", Value: nil})
		conditions = append(conditions, bson.E{Key: "is_publish", Value: false})
	case "bin":
		conditions = append(conditions, bson.E{Key: "deleted_at", Value: bson.M{"$ne": nil}})
	}

	return conditions
}

// buildPostAggregationPipeline 构建文章聚合管道
func (r *PostRepository) buildPostAggregationPipeline(cond bson.D, page *domain.PostQueryPage, skip, limit int64) mongo.Pipeline {
	sort := bson.D{{Key: "is_top", Value: -1}, {Key: "created_at", Value: -1}}
	if page.Field != "" {
		sort = append(sort, bson.E{Key: page.Field, Value: r.OrderConvertToInt(page.Order)})
	}

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

	if page.LabelName != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "type", Value: "tag"},
				{Key: "name", Value: page.LabelName},
			}}}},
		}}})
	}

	if page.CategoryName != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
			{Key: "category.type", Value: "category"},
			{Key: "category.name", Value: page.CategoryName},
		}}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$sort", Value: sort}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	return pipeline
}

// GetList 获取文章列表
func (r *PostRepository) GetList(ctx context.Context, page *domain.PostQueryPage, postType string) ([]*domain.PostDetail, int64, error) {
	// 构建查询条件
	cond := r.buildPostQueryCondition(page, postType)

	// 计算分页参数
	skip := (page.PageNo - 1) * page.PageSize
	limit := page.PageSize

	// 构建聚合管道
	pagePipeline := r.buildPostAggregationPipeline(cond, page, skip, limit)

	// 执行查询
	hasTagFilter := page.LabelName != ""
	hasCategoryFilter := page.CategoryName != ""
	posts, count, err := r.dao.GetListWithFilter(ctx, pagePipeline, cond, hasTagFilter, page.LabelName, hasCategoryFilter, page.CategoryName)
	if err != nil {
		return nil, 0, err
	}

	return r.PostCategoryTagsDOToPostDetailList(posts), count, nil
}

func (r *PostRepository) GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error) {
	cond := bson.D{
		{Key: "deleted_at", Value: nil},
		{Key: "is_publish", Value: true},
	}

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
		{{Key: "$sort", Value: bson.D{{Key: "is_top", Value: -1}, {Key: "created_at", Value: -1}}}},
	}

	posts, _, err := r.dao.GetList(ctx, pipeline, cond)
	if err != nil {
		return nil, err
	}
	return r.PostCategoryTagsDOToPostDetailList(posts), nil
}

func (r *PostRepository) FindByAlias(ctx context.Context, alias string) (*domain.Post, error) {
	post, err := r.dao.FindByAlias(ctx, alias)
	if err != nil {
		return nil, err
	}
	return r.PostDOToPostDomain(post), nil
}

func (r *PostRepository) PostDomainToPostDO(post *domain.Post) *dao.Post {
	return &dao.Post{
		CreatedAt:   post.CreatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryID:  post.CategoryId,
		TagsID:      post.TagsId,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDomainToUpdatePostDO(post *domain.Post) *dao.UpdatePost {
	return &dao.UpdatePost{
		CreatedAt:   post.CreatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryId:  post.CategoryId,
		TagsId:      post.TagsId,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostDOToPostDomainList(posts []*dao.Post) []*domain.Post {
	return lo.Map(posts, func(post *dao.Post, _ int) *domain.Post {
		return r.PostDOToPostDomain(post)
	})
}

func (r *PostRepository) PostDOToPostDomain(post *dao.Post) *domain.Post {
	return &domain.Post{
		Id:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Content:     post.Content,
		Description: post.Description,
		Author:      post.Author,
		Alias:       post.Alias,
		CategoryId:  post.CategoryID,
		TagsId:      post.TagsID,
		IsPublish:   post.IsPublish,
		IsTop:       post.IsTop,
		Thumbnail:   post.Thumbnail,
	}
}

func (r *PostRepository) PostCategoryTagsDOToPostDetail(postCategoryTags *dao.PostCategoryTags) *domain.PostDetail {
	return &domain.PostDetail{
		Id:          postCategoryTags.Id,
		CreatedAt:   postCategoryTags.CreatedAt,
		UpdatedAt:   postCategoryTags.UpdatedAt,
		Title:       postCategoryTags.Title,
		Content:     postCategoryTags.Content,
		Description: postCategoryTags.Description,
		Author:      postCategoryTags.Author,
		Alias:       postCategoryTags.Alias,
		Category:    postCategoryTags.Category,
		Tags:        postCategoryTags.Tags,
		Thumbnail:   postCategoryTags.Thumbnail,
		IsPublish:   postCategoryTags.IsPublish,
		IsTop:       postCategoryTags.IsTop,
	}
}

func (r *PostRepository) PostCategoryTagsDOToPostDetailList(postCategoryTags []*dao.PostCategoryTags) []*domain.PostDetail {
	return lo.Map(postCategoryTags, func(postCategoryTag *dao.PostCategoryTags, _ int) *domain.PostDetail {
		return r.PostCategoryTagsDOToPostDetail(postCategoryTag)
	})
}

func (r *PostRepository) OrderConvertToInt(order string) int {
	switch order {
	case "DESC":
		return -1
	case "ASC":
		return 1
	default:
		return -1
	}
}
