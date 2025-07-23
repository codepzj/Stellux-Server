package web

type CommentRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Url      string `json:"url"`
	Avatar   string `json:"avatar" binding:"required"`
	Content  string `json:"content" binding:"required"`
	PostId    string `json:"post_id" binding:"required,bson_id"`
	CommentId string `json:"comment_id" binding:"required,bson_id"`
}

type CommentEditRequest struct {
	Id       string `json:"id" binding:"required,bson_id"`
	Content  string `json:"content" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Url      string `json:"url"`
	Avatar   string `json:"avatar" binding:"required"`
}
