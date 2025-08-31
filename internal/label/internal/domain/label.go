package domain

import "go.mongodb.org/mongo-driver/v2/bson"

// Label 标签/分类实体
type Label struct {
	Id        bson.ObjectID // 标签ID
	LabelType string        // 标签类型: "category" | "tag"
	Name      string        // 标签名称
}

type LabelPostCount struct {
	Id        bson.ObjectID
	LabelType string
	Name      string
	Count     int
}
