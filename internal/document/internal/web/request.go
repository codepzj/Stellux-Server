package web

type DocumentRequest struct {
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content"`
	DocumentType string `json:"document_type" binding:"required"`
	ParentID     string `json:"parent_id"`
	DocumentID   string `json:"document_id"`
}

type DocumentRootRequest struct {
	Title        string `json:"title" binding:"required"`
	Alias        string `json:"alias" binding:"required"`
	Description  string `json:"description"`
	Thumbnail    string `json:"thumbnail"`
	DocumentType string `json:"document_type" binding:"required"`
	IsPublic     bool   `json:"is_public"`
}

type DocumentRootEditRequest struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Alias string `json:"alias" binding:"required"`
	Description string `json:"description"`
	Thumbnail string `json:"thumbnail"`
	DocumentType string `json:"document_type" binding:"required"`
	IsPublic bool `json:"is_public"`
}

type DocumentRootRenameRequest struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Alias string `json:"alias" binding:"required"`
	Description string `json:"description"`
	Thumbnail string `json:"thumbnail"`
	IsPublic bool `json:"is_public"`
}

type DeleteDocumentRequest struct {
	DocumentID string `json:"document_id" binding:"required"`
}

type DeleteDocumentListRequest struct {
	DocumentIDList []string `json:"document_id_list" binding:"required"`
}

type UpdateDocumentRequest struct {
	ID       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type RenameDocumentRequest struct {
	ID    string `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}
