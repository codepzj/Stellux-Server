package web

type DeleteFilesRequest struct {
	IDList []uint `json:"id_list" binding:"required"`
}
