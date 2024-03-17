package main

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/pupptmstr/RverStreaming/backend/pkg/api/v1"
	"log"
)

func main() {
	router := gin.Default()

	router.POST("/upload", v1.UploadFile)
	router.GET("/files", v1.ListFiles)
	router.GET("/file/:id", v1.GetFile)
	router.GET("/manifest/:id", v1.GetManifest)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("started on port :8080")
}
