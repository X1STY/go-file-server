package app

import (
	"github.com/gin-gonic/gin"
	"go-file-server/internal/app/file"
	"go-file-server/pkg/postgres"
)

func StartApplication() {
	db := postgres.SetupDatabase()
	fileHandler := file.NewFileHandler(db)

	server := SetupServer(fileHandler)
	server.Run(":8080")
}

func SetupServer(fileHandler *file.FileHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/files", fileHandler.GetAllFiles)
	router.POST("/file", fileHandler.UploadFile)
	router.GET("/file/:id", fileHandler.ReturnFile)

	return router
}
