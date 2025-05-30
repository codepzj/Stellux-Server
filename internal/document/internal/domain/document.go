package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Document struct {
	ID           bson.ObjectID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Title        string        // 文档标题
	Content      string        // 文档内容
	DocumentType string        // 文档类型
	ParentID     bson.ObjectID // 父文档id
	DocumentID   bson.ObjectID // 文档id
}

type DocumentRoot struct {
	ID           bson.ObjectID
	CreatedAt    time.Time // 创建时间
	UpdatedAt    time.Time // 更新时间
	Title        string    // 文档标题
	Alias        string    // 文档别名
	Description  string    // 文档描述
	DocumentType string    // 文档类型
	Thumbnail    string    // 文档缩略图
	IsPublic     bool      // 是否公开
}
