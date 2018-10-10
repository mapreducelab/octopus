package process

import (
	"octopus/config"
	"octopus/services/store"
	"os"

	"github.com/Shopify/sarama"
	"github.com/go-kit/kit/log"
)

// A Processor service processes messages from streaming message brokers.
type Processor interface {
	Process(con config.Connection, msg *sarama.ConsumerMessage) (err error)
}

type basicProcessorService struct {
	payload []byte
}

// Process messages from streaming services.
func (bps *basicProcessorService) Process(con config.Connection, msg *sarama.ConsumerMessage) (err error) {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	s := store.LogMiddleware{
		Logger: logger,
		Next:   store.NewStoreService(),
	}
	bps.payload = append(bps.payload, msg.Value...)

	if len(bps.payload) > con.UploadSize {
		s.Minio(con.Minio, bps.payload, msg.Topic)
		bps.payload = []byte{}
	}

	return err
}

// NewProcessorService processes and returns messages.
func NewProcessorService() Processor {
	return &basicProcessorService{}
}
