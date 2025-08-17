package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	Id          bson.ObjectID // 文档id
	CreatedAt   time.Time     // 创建时间
	UpdatedAt   time.Time     // 更新时间
	DeletedAt   time.Time     // 删除时间
	Title       string        // 文档标题
	Description string        // 文档描述
	Thumbnail   string        // 文档缩略图
	Alias       string        // 文档别名
	Sort        int           // 文档排序
	IsPublic    bool          // 是否公开
	IsDeleted   bool          // 是否删除
}

// Page 分页查询参数
type Page struct {
	PageNo   int64 `json:"pageNo"`   // 页码
	PageSize int64 `json:"pageSize"` // 每页大小
}
