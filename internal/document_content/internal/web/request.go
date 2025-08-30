package web

// 创建文档内容请求
type CreateDocumentContentRequest struct {
	DocumentId  string `json:"document_id" binding:"required,bson_id"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`     // 内容不是必填的
	Description string `json:"description"` // 描述不是必填的
	Alias       string `json:"alias"`
	ParentId    string `json:"parent_id" binding:"omitempty,bson_id"` // 允许为空，但如果有值必须是有效的ObjectID
	IsDir       bool   `json:"is_dir"`
	Sort        int    `json:"sort"`
}

// 更新文档内容请求
type UpdateDocumentContentRequest struct {
	Id          string `json:"id" binding:"required,bson_id"`
	DocumentId  string `json:"document_id" binding:"required,bson_id"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`     // 内容不是必填的
	Description string `json:"description"` // 描述不是必填的
	Alias       string `json:"alias"`
	ParentId    string `json:"parent_id" binding:"omitempty,bson_id"` // 允许为空，但如果有值必须是有效的ObjectID
	IsDir       bool   `json:"is_dir"`
	Sort        int    `json:"sort"`
}
