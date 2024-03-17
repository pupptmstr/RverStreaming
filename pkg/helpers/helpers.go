package helpers

import (
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"log"
	"os"
)

var FFMpegOutputBaseDirName = "ffmpeg-output"
var MinioClient *minio.Client
var bucketNameEnvVarName = "STORAGE_BUCKET_NAME"
var endpointEnvVarName = "STORAGE_ENDPOINT"
var accessKeyIDEnvVarName = "STORAGE_ACCESS_KEY_ID"
var secretAccessKeyEnvVarName = "STORAGE_SECRET_ACCESS_KEY"
var BucketName string
var endpoint string
var accessKeyID string
var secretAccessKey string

func init() {
	var exists bool
	BucketName, exists = os.LookupEnv(bucketNameEnvVarName)
	if !exists {
		log.Fatalf("Failed to get %s environment variable, no suck environment variable was found", bucketNameEnvVarName)
	}

	endpoint, exists = os.LookupEnv(endpointEnvVarName)
	if !exists {
		log.Fatalf("Failed to get %s environment variable, no suck environment variable was found", endpointEnvVarName)
	}

	accessKeyID, exists = os.LookupEnv(accessKeyIDEnvVarName)
	if !exists {
		log.Fatalf("Failed to get %s environment variable, no suck environment variable was found", accessKeyIDEnvVarName)
	}

	secretAccessKey, exists = os.LookupEnv(secretAccessKeyEnvVarName)
	if !exists {
		log.Fatalf("Failed to get %s environment variable, no suck environment variable was found", secretAccessKeyEnvVarName)
	}
	var err error
	MinioClient, err = InitMinioClient(BucketName, endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	EnsureDir(FFMpegOutputBaseDirName)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func EnsureDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0777)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dirName, err)
			return err
		}
	}

	return nil
}

func EnsureFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", filePath, err)
			return err
		}
		defer file.Close()
	}
	return nil
}
