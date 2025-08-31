package domain

import (
	"time"

	"github.com/codepzj/stellux/server/internal/label"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Post 文章
type Post struct {
	Id          bson.ObjectID   // 文章ID
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	Title       string          // 标题
	Content     string          // 内容
	Description string          // 描述
	Author      string          // 作者
	Alias       string          // 别名
	CategoryId  bson.ObjectID   // 分类ID
	TagsId      []bson.ObjectID // 标签ID
	IsPublish   bool            // 是否发布
	IsTop       bool            // 是否置顶
	Thumbnail   string          // 缩略图
}

type PostDetail struct {
	Id          bson.ObjectID   // 文章ID
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	Title       string          // 标题
	Content     string          // 内容
	Description string          // 描述
	Author      string          // 作者
	Alias       string          // 别名
	CategoryId  bson.ObjectID   // 分类Id
	Category    label.Domain    // 分类
	TagsId      []bson.ObjectID // 标签Id
	Tags        []label.Domain  // 标签
	IsPublish   bool            // 是否发布
	IsTop       bool            // 是否置顶
	Thumbnail   string          // 缩略图
}

type Page struct {
	PageNo    int64
	PageSize  int64
	Field     string
	Order     string
	Keyword   string
	LabelName string // 标签名称，用于过滤LabelType为"tag"的标签
}
