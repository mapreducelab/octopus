package store

import (
	"bytes"
	"fmt"
	"octopus/config"
	"time"

	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
)

// A Store service saves files to different storage systems.
type Store interface {
	Minio(config config.Minio, b []byte, topic string) (o string, err error)
}

type basicStoreService struct{}

// Minio stores files to Minio storage system.
func (bss *basicStoreService) Minio(config config.Minio, b []byte, topic string) (object string, err error) {
	minioClient, err := minio.New(config.Endpoint, config.AccessKeyID, config.SecretAccessKey, config.UseSSL)
	if err != nil {
		return "", errors.Wrap(err, "failed to initialize minio client")
	}
	timestamp := time.Now().UTC()
	object = fmt.Sprintf("nac/eng/kafka/%s/p04/%s-%v.json", topic, topic, timestamp.Format(time.RFC3339))

	r := bytes.NewReader(b)

	_, err = minioClient.PutObject(config.BucketName, object, r, int64(len(b)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return "", errors.Wrap(err, "failed to upload object Minio storage")
	}

	return object, err
}

// NewStoreService returns a naive, stateless implementation of StoreService.
func NewStoreService() Store {
	return &basicStoreService{}
}
