package storage

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/pupptmstr/RverStreaming/backend/pkg/fileutils"
	constants "github.com/pupptmstr/RverStreaming/backend/pkg/helpers"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// SaveFileToObjectStorageWithName saves file to object storage with given fileID in format "dirPath/<uuid>.<orig-file-ext>"
func SaveFileToObjectStorageWithName(objectStorage *minio.Client, file *multipart.FileHeader, fileID string, dirPath string) (*fileutils.FileData, error) {
	fileExt := filepath.Ext(file.Filename)
	var err error
	tempFile, err := fileutils.CreateTempFileFromUploadedFile(file)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	tempFilePath := tempFile.Name()
	var fullDirPath string
	if dirPath != "" {
		fullDirPath = dirPath + "/"
	} else {
		fullDirPath = ""
	}

	fullFileName := fullDirPath + fileID + fileExt
	_, err = objectStorage.FPutObject(context.Background(), constants.BucketName, fullFileName, tempFilePath, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println(fullFileName + " err: Failed to upload file to data storage")
		return nil, err
	}

	defer os.Remove(tempFilePath)

	fileData := &fileutils.FileData{
		FileName:              fileID,
		OriginalFileSizeBytes: int(file.Size),
		AvailableQualities:    []string{},
	}
	return fileData, nil
}

// SaveFileToObjectStorage saves the file to the object storage with generated UUID in format "<uuid>.<orig-file-ext>"
func SaveFileToObjectStorage(objectStorage *minio.Client, file *multipart.FileHeader) (*fileutils.FileData, error) {
	fileID := constants.GenerateUUID()
	return SaveFileToObjectStorageWithName(objectStorage, file, fileID, "")
}

func GetOriginalFileDataById(objectStorage *minio.Client, fileID string) (*minio.ObjectInfo, error) {
	fileInfo, err := objectStorage.StatObject(context.Background(), constants.BucketName, fileID, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &fileInfo, nil
}

func GetOriginalFileById(objectStorage *minio.Client, fileID string) (*minio.Object, error) {
	file, err := objectStorage.GetObject(context.Background(), constants.BucketName, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetFileIDsList(objectStorage *minio.Client) []fileutils.FileData {
	files := objectStorage.ListObjects(context.Background(), constants.BucketName, minio.ListObjectsOptions{})
	var fileList []fileutils.FileData
	for object := range files {
		if strings.HasSuffix(object.Key, "/") {
			continue
		}

		id := strings.TrimSuffix(object.Key, filepath.Ext(object.Key))
		manifestPath := id + "/" + "manifest.mpd"
		bitrates, err := ReadDashManifest(objectStorage, manifestPath)
		if err != nil {
			fmt.Printf("Error while reading file manifest with id: %s; Err: %s", id, err)
			continue
		}
		fileData := fileutils.FileData{
			FileName:              id,
			OriginalFileSizeBytes: int(object.Size),
			AvailableQualities:    bitrates,
		}
		fileList = append(fileList, fileData)
	}
	return fileList
}

func ReadDashManifest(objectStorage *minio.Client, manifestPath string) ([]string, error) {
	reader, err := objectStorage.GetObject(context.Background(), constants.BucketName, manifestPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var mpd fileutils.MPD
	if err := xml.NewDecoder(reader).Decode(&mpd); err != nil {
		return nil, err
	}

	var qualities []string
	for _, period := range mpd.Periods {
		for _, adaptationSet := range period.AdaptationSets {
			for _, representation := range adaptationSet.Representations {
				if representation.Width == "" || representation.Height == "" {
					//skip
					continue
				}
				quality := representation.Height + "p"
				if representation.FrameRate != "" {
					quality += representation.FrameRate
				}
				qualities = append(qualities, quality)
			}
		}
	}

	return qualities, nil
}

func GetManifest(objectStorage *minio.Client, fileId string) (*minio.Object, error) {
	manifestFilePath := fileId + "/" + "manifest.mpd"
	manifestFile, err := objectStorage.GetObject(context.Background(), constants.BucketName, manifestFilePath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return manifestFile, nil
}
