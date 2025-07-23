package web

type CommentRequest struct {
	Content   string `json:"content" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Url       string `json:"url" binding:"required"`
	Avatar    string `json:"avatar" binding:"required"`
	PostId    string `json:"post_id" binding:"required,bson_id"`
	CommentId string `json:"comment_id" binding:"required,bson_id"`
}

type CommentEditRequest struct {
	Id       string `json:"id" binding:"required,bson_id"`
	Content  string `json:"content" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Url      string `json:"url" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}
