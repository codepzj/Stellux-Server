package web

type DocumentCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	Alias       string `json:"alias" binding:"required"`
	Sort        int    `json:"sort" binding:"required,gt=0"`
	IsPublic    bool   `json:"is_public"`
}

type DocumentUpdateRequest struct {
	Id          string `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Thumbnail   string `json:"thumbnail"`
	Alias       string `json:"alias" binding:"required"`
	Sort        int    `json:"sort" binding:"required,gt=0"`
	IsPublic    bool   `json:"is_public"`
}
