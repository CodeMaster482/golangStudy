package handler

import (
	"context"
	"fmt"
	mClient "main/internal/minio"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (h *Handler) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(header.Filename)
	if _, err := h.MinioCli.PutObject(ctx, mClient.BucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", mClient.MinioHost, mClient.BucketName, objectName), nil
}

func (h *Handler) DeleteImage(ctx context.Context, objectName string) error {
	return h.MinioCli.RemoveObject(ctx, mClient.BucketName, objectName, minio.RemoveObjectOptions{})
}
