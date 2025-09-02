package repository

import (
	"context"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/bsonx"
	"github.com/chenmingyong0423/go-mongox/v2/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
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
	GetList(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error)
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
func (r *PostRepository) buildPostQueryCondition(page *domain.Page, postType string) bson.D {
	builder := query.NewBuilder().
		Or(
			query.RegexOptions("title", page.Keyword, "i"),
			query.RegexOptions("description", page.Keyword, "i"),
		)

	switch postType {
	case "publish":
		builder = builder.And(query.Eq("deleted_at", nil)).And(query.Eq("is_publish", true))
	case "draft":
		builder = builder.And(query.Eq("deleted_at", nil)).And(query.Eq("is_publish", false))
	case "bin":
		builder = builder.And(query.Ne("deleted_at", nil))
	}

	return builder.Build()
}

// buildPostAggregationPipeline 构建文章聚合管道
func (r *PostRepository) buildPostAggregationPipeline(cond bson.D, page *domain.Page, skip, limit int64) mongo.Pipeline {
	sortBuilder := bsonx.NewD().Add("is_top", -1).Add("created_at", -1)
	if page.Field != "" {
		sortBuilder.Add(page.Field, r.OrderConvertToInt(page.Order))
	}

	stageBuilder := aggregation.NewStageBuilder().Match(cond).
		Lookup("label", "category", &aggregation.LookUpOptions{
			LocalField:   "category_id",
			ForeignField: "_id",
		}).
		Unwind("$category", &aggregation.UnWindOptions{
			PreserveNullAndEmptyArrays: true,
		}).
		Lookup("label", "tags", &aggregation.LookUpOptions{
			LocalField:   "tags_id",
			ForeignField: "_id",
		})

	// 添加标签过滤条件
	if page.LabelName != "" {
		stageBuilder = stageBuilder.Match(query.ElemMatch("tags", query.And(
			query.Eq("type", "tag"),
			query.Eq("name", page.LabelName),
		)))
	}

	return stageBuilder.Sort(sortBuilder.Build()).Skip(skip).Limit(limit).Build()
}

// GetList 获取文章列表
func (r *PostRepository) GetList(ctx context.Context, page *domain.Page, postType string) ([]*domain.PostDetail, int64, error) {
	// 构建查询条件
	cond := r.buildPostQueryCondition(page, postType)

	// 计算分页参数
	skip := (page.PageNo - 1) * page.PageSize
	limit := page.PageSize

	// 构建聚合管道
	pagePipeline := r.buildPostAggregationPipeline(cond, page, skip, limit)

	// 执行查询
	hasTagFilter := page.LabelName != ""
	posts, count, err := r.dao.GetListWithTagFilter(ctx, pagePipeline, cond, hasTagFilter, page.LabelName)
	if err != nil {
		return nil, 0, err
	}

	return r.PostCategoryTagsDOToPostDetailList(posts), count, nil
}

func (r *PostRepository) GetAllPublishPost(ctx context.Context) ([]*domain.PostDetail, error) {
	var cond bson.D
	builder := query.NewBuilder().And(query.Eq("deleted_at", nil)).And(query.Eq("is_publish", true))
	cond = builder.Build()
	sortBuilder := bsonx.NewD().Add("is_top", -1).Add("created_at", -1)
	pipeline := aggregation.NewStageBuilder().Match(cond).
		Lookup("label", "category", &aggregation.LookUpOptions{
			LocalField:   "category_id",
			ForeignField: "_id",
		}).
		Unwind("$category", &aggregation.UnWindOptions{
			PreserveNullAndEmptyArrays: true,
		}).
		Lookup("label", "tags", &aggregation.LookUpOptions{
			LocalField:   "tags_id",
			ForeignField: "_id",
		}).
		Sort(sortBuilder.Build()).Build()
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

// PostDomain2PostDO 将domain.Post转换为dao.Post
func (r *PostRepository) PostDomainToPostDO(post *domain.Post) *dao.Post {
	return &dao.Post{
		Model:       mongox.Model{CreatedAt: post.CreatedAt},
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

// PostDOToPostDomain 将dao.Post转换为domain.Post
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

// PostCategoryTagsDOToPostDetail 将dao.PostCategoryTags转换为domain.PostDetail
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
