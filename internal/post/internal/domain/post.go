package domain

import (
	"time"

	"github.com/codepzj/Stellux-Server/internal/label"
	"github.com/codepzj/Stellux-Server/internal/pkg/apiwrap"
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

// PostQueryPage Post模块的分页查询参数（包含特殊过滤字段）
type PostQueryPage struct {
	apiwrap.Page        // 嵌入通用分页参数
	LabelName    string // 标签名称过滤
	CategoryName string // 分类名称过滤
}
