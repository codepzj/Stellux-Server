package apiwrap

type Page struct {
	// 当前页
	PageNo int64 `form:"page_no" json:"page_no" binding:"required,gte=1"`
	// 每页条数
	PageSize int64 `form:"page_size" json:"page_size" binding:"required,gte=1"`
	// 排序字段
	Field string `form:"field" json:"field,omitempty" `
	// 排序方式
	Order string `form:"order" json:"order,omitempty" binding:"omitempty,oneof=ASC DESC"`
	// 搜索内容
	Keyword string `form:"keyword" json:"keyword,omitempty" `
}

type PageVO[T any] struct {
	Page
	// 总条数
	TotalCount int64 `json:"total_count"`
	// 总页数
	TotalPage int64 `json:"total_page"`
	List      []*T  `json:"list"`
}

func ToPageVO[T any](pageNo int64, pageSize int64, totalCount int64, list []*T) *PageVO[T] {
	return &PageVO[T]{
		Page: Page{
			PageNo:   pageNo,
			PageSize: pageSize,
		},
		TotalCount: totalCount,
		TotalPage:  (totalCount + pageSize - 1) / pageSize,
		List:       list,
	}
}
