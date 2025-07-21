package web

type DocumentCreateDto struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Alias       string `json:"alias" binding:"required"`
	Sort        int    `json:"sort" binding:"required,gt=0"`
	IsPublic    bool   `json:"is_public"`
}

type DocumentUpdateDto struct {
	Id          string `json:"id" binding:"required,bsonId"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Alias       string `json:"alias" binding:"required"`
	Sort        int    `json:"sort" binding:"required,gt=0"`
	IsPublic    bool   `json:"is_public"`
}
