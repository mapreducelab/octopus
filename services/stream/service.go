package stream

import (
	"fmt"
	stdlog "log"
	"octopus/config"
	"octopus/services/process"
	"os"
	"os/signal"

	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-kit/kit/log"
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

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.ClientID = con.KafkaCon.ClientID

	topics := []string{con.KafkaCon.Topic}
	consumerGroup := fmt.Sprintf("%s-consumer-group", con.KafkaCon.ClientID)

	consumer, err := cluster.NewConsumer(con.KafkaCon.Brokers, consumerGroup, topics, config)
	if err != nil {
		return errors.Wrap(err, "failed to create consumer")
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			stdlog.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			stdlog.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	pr := process.LogMiddleware{
		Logger: logger,
		Next:   process.NewProcessorService(),
	}

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				if err := pr.Process(con, msg); err != nil {
					stdlog.Printf("Processing failed :%+v", err)
				}
				consumer.MarkOffset(msg, "")
			}
		case <-signals:
			return
		}
	}
}

// NewStreamerService returns a naive, stateless implementation of StreamerService.
func NewStreamerService() Streamer {
	return &streamerService{}
}
