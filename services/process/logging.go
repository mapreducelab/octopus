package process

import (
	"octopus/config"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/go-kit/kit/log"
)

// LogMiddleware describes a service middleware.
type LogMiddleware struct {
	Logger log.Logger
	Next   Processor
}

// Process messages from streaming services.
func (lm LogMiddleware) Process(con config.Connection, msg *sarama.ConsumerMessage) (err error) {
	defer func(begin time.Time) {
		method := "store.Process"
		logger := log.With(lm.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
		_ = logger.Log(
			"method", method,
			"Topic", msg.Topic,
			"Partition", msg.Partition,
			"size", len(msg.Value),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = lm.Next.Process(con, msg)
	return
}
