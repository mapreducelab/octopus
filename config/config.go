package config

import (
	"io/ioutil"

	"github.com/pkg/errors"

	yaml "gopkg.in/yaml.v2"
)

// A Connection is a generic configuration object for establishing
// connection to streaming service.
type Connection struct {
	// StreamingService is a type of streaming service, foe example "KAFKA".
	StreamingService string `yaml:"streamingService"`

	// Apache Kafka configuration object.
	KafkaCon KafkaCon `yaml:"kafkaCon"`

	// Configuration for Minio object store.
	Minio Minio `yaml:"minio"`

	// Limit size to upload to the storage system.
	UploadSize int `yaml:"uploadSize"`
}

// A KafkaCon represents a configuration object to establish
// a connection to Apache Kafka.
type KafkaCon struct {
	// Apache Kafka brokers to establish a connection.
	Brokers []string `yaml:"brokers"`

	// Kafka topic to connect to.
	Topic string `yaml:"topic"`

	// Kafka partision to connect to.
	Partision int32 `yaml:"partision"`

	// Execution mode, if "false" will be lunched only one instance of Streamer
	// per Kafka topic, otherwise one instance of Streamer per topic's partision.
	Distributed bool `yaml:"distributed"`
}

// A Minio represents a configuration object to establish
// a connection to Minio s3 storage.
type Minio struct {
	// Minio url to connect to server.
	Endpoint string `yaml:"endpoint"`

	// Access Key from Minio.
	AccessKeyID string `yaml:"accessKeyID"`

	// Secret Access Key from Minio.
	SecretAccessKey string `yaml:"secretAccessKey"`

	// Bucket name.
	BucketName string `yaml:"bucketName"`

	// Use SSL or not.
	UseSSL bool `yaml:"UseSSL"`
}

// GetCon unmarshal yaml configuration file into struct.
func (c *Connection) GetCon(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read configuration yaml file: %s", path)
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		return errors.Wrapf(err, "failed to unmarshal conf file")
	}
	return nil
}
