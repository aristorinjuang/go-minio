package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	err    error
	useSSL bool
)

func init() {
	if err = godotenv.Load(); err != nil {
		log.Panic(err)
	}

	if useSSL, err = strconv.ParseBool(os.Getenv("MINIO_SSL")); err != nil {
		log.Panic(err)
	}
}

func main() {
	ctx := context.Background()
	endpoint := os.Getenv("MINIO_HOST")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Panic(err)
	}

	bucket := os.Getenv("MINIO_BUCKET")
	if err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		if exists, errBucketExists := minioClient.BucketExists(ctx, bucket); errBucketExists == nil && exists {
			log.Println("BUCKET EXISTS:", bucket)
		} else {
			log.Panic(err)
		}
	}
	log.Println("BUCKET CREATED:", bucket)

	fileName := os.Getenv("FILENAME")
	filePath := os.Getenv("FILEPATH")
	mimeType := os.Getenv("MIMETYPE")

	uploadInfo, err := minioClient.FPutObject(ctx, bucket, fileName, filePath, minio.PutObjectOptions{
		ContentType: mimeType,
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println("OBJECT CREATED:", uploadInfo)
}
