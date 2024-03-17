package fileutils

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v7"
	"github.com/pupptmstr/RverStreaming/backend/pkg/helpers"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
)

type FileData struct {
	FileName              string   `json:"file_name"`
	OriginalFileSizeBytes int      `json:"original_file_size_bytes"`
	AvailableQualities    []string `json:"available_qualities"`
}

type MPD struct {
	XMLName xml.Name `xml:"MPD"`
	Periods []Period `xml:"Period"`
}

type Period struct {
	AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
	Representations []Representation `xml:"Representation"`
}

type Representation struct {
	Bandwidth string `xml:"bandwidth,attr"`
	Width     string `xml:"width,attr"`
	Height    string `xml:"height,attr"`
	FrameRate string `xml:"frameRate,attr"`
}

func IsFileVideoOrAudio(file *multipart.FileHeader) (bool, error) {
	openedFile, err := file.Open()
	if err != nil {
		fmt.Println("Failed to open the file")
		return false, err
	}
	defer openedFile.Close()

	buffer := make([]byte, 261)
	if _, err := openedFile.Read(buffer); err != nil {
		return false, err
	}

	kind, _ := filetype.Match(buffer)
	if kind == filetype.Unknown {
		return false, nil
	}
	res := kind.MIME.Type == "video" || kind.MIME.Type == "audio"
	return res, nil
}

func TranscodeAndSegmentFile(fileData FileData, file *multipart.FileHeader) error {
	outputDirName := helpers.FFMpegOutputBaseDirName + "/" + fileData.FileName
	err := helpers.EnsureDir(outputDirName)
	if err != nil {
		return err
	}

	tempFile, err := CreateTempFileFromUploadedFile(file)
	if err != nil {
		return err
	}

	defer tempFile.Close()

	tempFilePath := tempFile.Name()
	manifestFilePath := outputDirName + "/manifest.mpd"
	cmd := exec.Command(
		"ffmpeg",
		"-i", tempFilePath,
		"-map", "0:v", "-map", "0:a",
		"-b:v:0", "1500k", "-filter:v:0", "scale=-1:720",
		"-b:v:1", "600k", "-filter:v:1", "scale=-1:360",
		"-var_stream_map", "v:0,a:0 v:1,a:1",
		"-f", "dash", manifestFilePath,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("FFmpeg error:", err)
		fmt.Println("FFmpeg stderr:", stderr.String())
		return err
	}

	files, err := os.ReadDir(outputDirName)
	if err != nil {
		return err
	}

	for _, dirEntry := range files {
		if !dirEntry.IsDir() {
			localFilePath := outputDirName + "/" + dirEntry.Name()
			objectName := fileData.FileName + "/" + dirEntry.Name()

			_, err = helpers.MinioClient.FPutObject(context.Background(), helpers.BucketName, objectName, localFilePath, minio.PutObjectOptions{})
			if err != nil {
				fmt.Println(objectName + " err: Failed to upload file to data storage")
				return err
			}
		}
	}

	err = os.RemoveAll(outputDirName)
	if err != nil {
		fmt.Println("Failed to remove ffmpeg created files:", err)
	}

	return nil
}

func CreateTempFileFromUploadedFile(file *multipart.FileHeader) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "upload-*"+file.Filename)
	if err != nil {
		fmt.Println(file.Filename + " err: Failed to create a temp file")
		return nil, err
	}

	uploadedFile, err := file.Open()
	if err != nil {
		fmt.Println(file.Filename + " err: Failed to open the uploaded file")
		tempFile.Close()
		uploadedFile.Close()
		return nil, err
	}

	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		fmt.Println(file.Filename + " err: Failed to copy the uploaded file to the temp file")
		tempFile.Close()
		uploadedFile.Close()
		return nil, err
	}

	return tempFile, nil
}
