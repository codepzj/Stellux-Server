package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DocumentContent struct {
	Id          bson.ObjectID // 文档内容ID
	CreatedAt   time.Time     // 创建时间
	UpdatedAt   time.Time     // 更新时间
	DeletedAt   time.Time     // 删除时间
	DocumentId  bson.ObjectID // 文档ID
	Title       string        // 文档标题
	Content     string        // 文档内容
	Description string        // 文档描述
	Alias       string        // 文档别名
	ParentId    bson.ObjectID // 父级ID
	IsDir       bool          // 是否是目录
	Sort        int           // 排序
	IsDeleted   bool          // 是否删除
}

// Page 分页查询参数
type Page struct {
	PageNo   int64 `json:"page_no"`   // 页码
	PageSize int64 `json:"page_size"` // 每页大小
}
