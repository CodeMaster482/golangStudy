package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go"

	minioCli "main/internal/minio"
)

func (h *Handler) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(header.Filename)
	if _, err := h.MinioCli.PutObject(minioCli.BucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", minioCli.MinioHost, minioCli.BucketName, objectName), nil
}

func (h *Handler) DeleteImage(objectName string) error {
	return h.MinioCli.RemoveObject(minioCli.BucketName, objectName)
}
