package web

import "mime/multipart"

type UploadFilesRequest struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}

type DeleteFilesRequest struct {
	IDList []string `json:"id_list" binding:"required"`
}
