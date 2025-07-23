package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	Id        bson.ObjectID
	CreatedAt time.Time
	UpdatedAt time.Time
	Nickname  string        // 昵称
	Email     string        // 邮箱
	Avatar    string        // 头像
	Url       string        // 网站
	Content   string        // 内容
	PostId    bson.ObjectID // 帖子id
	CommentId bson.ObjectID // 回复的评论id
}
