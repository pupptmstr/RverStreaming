package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pupptmstr/RverStreaming/backend/pkg/fileutils"
	constants "github.com/pupptmstr/RverStreaming/backend/pkg/helpers"
	"github.com/pupptmstr/RverStreaming/backend/pkg/storage"
	"io"
	"net/http"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	isMediaFile, err := fileutils.IsFileVideoOrAudio(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while checking file format. Err: " + err.Error()})
		return
	}

	if !isMediaFile {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while checking file type, it's not media file"})
		return
	}

	uploadedFile, err := storage.SaveFileToObjectStorage(constants.MinioClient, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while uploading file to object storage. Err: " + err.Error()})
		return
	}

	err = fileutils.TranscodeAndSegmentFile(*uploadedFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while transcoding and segmenting file. Err: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "id": uploadedFile.FileName})
}

func ListFiles(c *gin.Context) {
	fileList := storage.GetFileIDsList(constants.MinioClient)
	htmlContent := "<html><body><h1>Available Files</h1><ul>"

	for _, file := range fileList {
		videoPageURL := fmt.Sprintf("/file/%s", file.FileName) // Путь к странице с видео
		htmlContent += fmt.Sprintf("<li><a href='%s'>%s</a></li>", videoPageURL, file.FileName)
	}

	htmlContent += "</ul></body></html>"
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent)
}

func GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is missing"})
		return
	}

	manifestURL := fmt.Sprintf("/manifest/%s", fileID)
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		    <link href="https://vjs.zencdn.net/8.10.0/video-js.css" rel="stylesheet" />
		</head>
		<body>
		    <video
		        id="videoPlayer"
		        class="video-js"
		        controls
		        preload="auto"
		        width="640"
		        height="264"
		        data-setup="{}">
		        <source src="%s" type="application/dash+xml">
		    </video>
		    <script>
		        var player = videojs('videoPlayer');
		    </script>
		</body>
		</html>`, manifestURL)

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, htmlContent) //TODO не работает, разобраться
}

func GetManifest(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is missing"})
		return
	}

	file, err := storage.GetManifest(constants.MinioClient, fileID)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get manifest file"})
		return
	}
	defer file.Close()

	c.Header("Content-Type", "application/dash+xml")
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return
	}
}
