package req

type CreatePostReq struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Alias       string `json:"alias"`
	CategoryID  int64  `json:"category_id"`
	IsPublish   int32  `json:"is_publish"`
	IsTop       int32  `json:"is_top"`
	Thumbnail   string `json:"thumbnail"`
}

type UpdatePostReq struct {
	ID int64 `json:"id"`
	CreatePostReq
}
