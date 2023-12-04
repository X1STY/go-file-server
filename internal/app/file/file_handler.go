package file

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type FileHandler struct {
	service *FileService
}

func NewFileHandler(db *sql.DB) *FileHandler {
	return &FileHandler{
		service: NewFileService(db),
	}
}

func (h *FileHandler) GetAllFiles(ctx *gin.Context) {
	files, err := h.service.GetAllFiles(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get files"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Files": files})
}

func (h *FileHandler) UploadFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad file request"})
		return
	}
	defer file.Close()

	err = h.service.UploadFile(header)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File upload successfully"})
}

func (h *FileHandler) ReturnFile(ctx *gin.Context) {
	fileIDStr := ctx.Param("id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	fileName, filePath, err := h.service.ReturnFile(fileID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.FileAttachment(filePath, fileName)
}
