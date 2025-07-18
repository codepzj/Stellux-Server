package web

// 创建文档内容请求
type CreateDocumentContentDto struct {
	DocumentId   string `json:"documentId" binding:"required,bsonId"`
	Title        string         `json:"title" binding:"required"`
	Content      string         `json:"content" binding:"required"`
	Description  string         `json:"description" binding:"required"`
	Version      string         `json:"version" binding:"required,version"`
	Alias        string         `json:"alias" binding:"required"`
	ParentId     string `json:"parentId" binding:"required,bsonId"`
	IsDir        bool           `json:"isDir"`
	Sort         int            `json:"sort"`
}

// 更新文档内容请求
type UpdateDocumentContentDto struct {
	Id           string `json:"id" binding:"required,bsonId"`
	DocumentId   string `json:"documentId" binding:"required,bsonId"`
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Description  string `json:"description" binding:"required"`
	Version      string `json:"version" binding:"required,version"`
	Alias        string `json:"alias" binding:"required"`
	ParentId     string `json:"parentId" binding:"required,bsonId"`
	IsDir        bool   `json:"isDir"`
	Sort         int    `json:"sort"`
}