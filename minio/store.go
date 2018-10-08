package minio

import (
	"bytes"
	"fmt"
	"log"
	"octopus/config"
	"time"

	"github.com/pkg/errors"

	"github.com/minio/minio-go"
)

// SaveToMinio accepts bytes abd uploads object to the Minio storage system.
func SaveToMinio(config config.Minio, b []byte, topic string) error {
	// Initialize minio client object.
	minioClient, err := minio.New(config.Endpoint, config.AccessKeyID, config.SecretAccessKey, config.UseSSL)
	if err != nil {
		return errors.Wrap(err, "failed to initialize minio client")
	}
	timestamp := time.Now().UTC()
	objectName := fmt.Sprintf("nac/eng/kafka/%s/p04/%s-%v.json", topic, topic, timestamp.Format(time.RFC3339))

	r := bytes.NewReader(b)

	n, err := minioClient.PutObject(config.BucketName, objectName, r, int64(len(b)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return errors.Wrap(err, "failed to upload object Minio storage")
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

	return nil
}
