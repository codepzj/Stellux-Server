package web

type DeleteFilesRequest struct {
	IDList []string `json:"id_list" binding:"required"`
}
