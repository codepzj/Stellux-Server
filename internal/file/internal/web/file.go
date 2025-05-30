package web

import (
	"fmt"

	"github.com/codepzj/stellux/server/internal/file/internal/service"
	"github.com/codepzj/stellux/server/internal/pkg/apiwrap"
	"github.com/gin-gonic/gin"
)

func NewFileHandler(serv service.IFileService) *FileHandler {
	return &FileHandler{
		serv: serv,
	}
}

type FileHandler struct {
	serv service.IFileService
}

func (h *FileHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.Static("/images", "./static/images")
	fileGroup := engine.Group("/file")
	{
		fileGroup.GET("/list", apiwrap.WrapWithBody(h.QueryFileList))
	}
	fileAdminGroup := engine.Group("/admin-api/file")
	{
		fileAdminGroup.POST("/upload", apiwrap.WrapWithBody(h.UploadFiles))
		fileAdminGroup.DELETE("/delete", apiwrap.WrapWithBody(h.DeleteFiles))
	}
}

func (h *FileHandler) UploadFiles(c *gin.Context, uploadFilesRequest *UploadFilesRequest) *apiwrap.Response[any] {
	files := uploadFilesRequest.Files
	if len(files) == 0 {
		return apiwrap.FailWithMsg(apiwrap.RuquestBadRequest, "未找到上传的文件")
	}
	err := h.serv.UploadFiles(c, files)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("文件上传成功")
}

func (h *FileHandler) QueryFileList(c *gin.Context, page *apiwrap.Page) *apiwrap.Response[any] {
	files, count, err := h.serv.QueryFileList(c, page)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithDetail[any](apiwrap.ToPageVO(page.PageNo, page.PageSize, count, h.FileDomainToVOList(files)), "文件列表查询成功")
}

func (h *FileHandler) DeleteFiles(c *gin.Context, deleteFilesRequest *DeleteFilesRequest) *apiwrap.Response[any] {
	idList := deleteFilesRequest.IDList
	fmt.Println(idList)
	err := h.serv.DeleteFiles(c, deleteFilesRequest.IDList)
	if err != nil {
		return apiwrap.FailWithMsg(apiwrap.RuquestInternalServerError, err.Error())
	}
	return apiwrap.SuccessWithMsg("文件删除成功")
}
