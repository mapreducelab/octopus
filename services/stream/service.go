package stream

import (
	"fmt"
	"log"
	"octopus/config"
	"octopus/minio"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

// A Streamer is a service to process streaming data from a different
// streaming services, such as Apache Kafka, RabbitMQ, NATS, etc.
type Streamer interface {
	Process(con config.Connection) (err error)
}

type streamerService struct{}

// A Process represents a mechanism to process streaming data.
func (c *streamerService) Process(con config.Connection) (err error) {
	if con.StreamingService != "KAFKA" {
		return errors.New(con.StreamingService + " is not supported yet by Streamer service")
	}
	if con.KafkaCon.Distributed == true {
		return errors.New("distributed mode is not yet supported")
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = "datalake"

	// Create new consumer
	master, err := sarama.NewConsumer(con.KafkaCon.Brokers, config)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumer, sErrors, err := consume(con.KafkaCon.Topic, master)
	if err != nil {
		return errors.Wrapf(err, "failed to consume kafka topic: %v", con.KafkaCon.Topic)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Get signal for finish
	doneCh := make(chan struct{})
	go func(uploadSize int) {
		var jsonBytes []byte
		for {
			select {
			case msg := <-consumer:
				jsonBytes = append(jsonBytes, msg.Value...)
				if len(jsonBytes) > uploadSize {
					minio.SaveToMinio(con.Minio, jsonBytes, msg.Topic)
					jsonBytes = []byte{}
				}
			case consumerError := <-sErrors:
				fmt.Println("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}(con.UploadSize)
	<-doneCh

	return nil
}

// A generator pattern, creates channels for each partition and starts to push messages.
func consume(topic string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError, error) {
	consumers := make(chan *sarama.ConsumerMessage)
	sErrors := make(chan *sarama.ConsumerError)

	partitions, err := master.Partitions(topic)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get partisions for topic: %v", topic)
	}

	for part := range partitions {
		consumer, err := master.ConsumePartition(topic, int32(part), sarama.OffsetNewest)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to consume partition: %v, of topic: %v", topic, int32(part))
		}
		log.Printf("Start consuming partision: %d of the topic: %v\n", int32(part), topic)
		go func(topic string, consumer sarama.PartitionConsumer, part int32) {
			for {
				select {
				case consumerError := <-consumer.Errors():
					sErrors <- consumerError
				case msg := <-consumer.Messages():
					consumers <- msg
				}
			}
		}(topic, consumer, int32(part))
	}

	return consumers, sErrors, nil
}

// NewStreamerService returns a naive, stateless implementation of StreamerService.
func NewStreamerService() Streamer {
	return &streamerService{}
}
