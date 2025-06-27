package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Friend struct {
	ID          bson.ObjectID
	Name        string // 网站名称
	Description string // 网站描述
	SiteUrl     string // 网站地址
	WebsiteType string // 网站类型
	AvatarUrl   string // 头像地址
	IsActive    bool   // 是否激活
}
