package web

type LabelRequest struct {
	ID        string `json:"id,omitempty"`
	LabelType string `json:"label_type" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

type Page struct {
	PageNo    int64  `form:"page_no" binding:"required,min=1"`
	PageSize  int64  `form:"page_size" binding:"required,min=1"`
	LabelType string `form:"label_type" binding:"required"`
}
