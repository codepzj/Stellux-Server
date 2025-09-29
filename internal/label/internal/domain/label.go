package domain

// Label 标签/分类实体
type Label struct {
	ID        uint   // 标签ID
	LabelType string // 标签类型: "category" | "tag"
	Name      string // 标签名称
}

type LabelPostCount struct {
	ID        uint
	LabelType string
	Name      string
	Count     int
}
