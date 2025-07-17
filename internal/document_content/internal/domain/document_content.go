package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DocumentContent struct {
	Id           bson.ObjectID  // 文档内容ID
	CreatedAt    time.Time     // 创建时间
	UpdatedAt    time.Time     // 更新时间
	DeletedAt    time.Time     // 删除时间
	DocumentId   bson.ObjectID // 文档ID
	Title        string         // 文档标题
	Content      string        // 文档内容
	Version      int       // 文档版本
	Alias        string       // 文档别名
	ParentID     bson.ObjectID // 父级ID
	IsDir        bool           // 是否是目录
	LikeCount    int            // 点赞数
	DislikeCount int            // 反对数
	CommentCount int            // 评论数
}

