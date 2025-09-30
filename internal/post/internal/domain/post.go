package domain

import (
	"time"

	"github.com/codepzj/Stellux-Server/internal/label"
)

// Post 文章
type Post struct {
	ID          uint            // 文章ID
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	Title       string          // 标题
	Content     string          // 内容
	Description string          // 描述
	Author      string          // 作者
	Alias       string          // 别名
	CategoryID  uint            // 分类ID
	TagsID      []uint          // 标签ID
	IsPublish   bool            // 是否发布
	IsTop       bool            // 是否置顶
	Thumbnail   string          // 缩略图
}

type PostDetail struct {
	ID          uint            // 文章ID
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 更新时间
	Title       string          // 标题
	Content     string          // 内容
	Description string          // 描述
	Author      string          // 作者
	Alias       string          // 别名
	CategoryID  uint            // 分类Id
	Category    label.Domain    // 分类
	TagsID      []uint          // 标签Id
	Tags        []label.Domain  // 标签
	IsPublish   bool            // 是否发布
	IsTop       bool            // 是否置顶
	Thumbnail   string          // 缩略图
}

// Page 分页查询参数
type Page struct {
	PageNo    int64  // 页码，从1开始
	PageSize  int64  // 每页大小
	Field     string // 排序字段
	Order     string // 排序方式：ASC/DESC
	Keyword   string // 搜索关键词
	LabelName string // 标签名称，用于过滤LabelType为"tag"的标签
}
