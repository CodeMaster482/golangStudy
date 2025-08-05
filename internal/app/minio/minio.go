package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

const (
	BucketName = "banknote-name-server"
	MinioHost  = "127.0.0.1:9000"
)

func NewMinioClient(logger *logrus.Logger) *minio.Client {
	accessKeyID := "minio"
	secretAccessKey := "minio124"

	// Create an Options struct instance
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	}

	minioClient, err := minio.New(MinioHost, opts)

	if err != nil {
		logger.Fatalf("error: %s", err)
	}
	location := "us-east-1"

	ctx := context.Background()

	err = minioClient.MakeBucket(ctx, BucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, err2 := minioClient.BucketExists(ctx, BucketName)
		if err2 == nil && exists {
			logger.Infof("We already own %s", BucketName)
		} else {
			logger.Fatalf("error: %s", err2)
		}
	} else {
		logger.Infof("Successfully created %s\n", BucketName)
	}

	return minioClient
}
