package stream

import (
	"octopus/config"
	"time"

	log "github.com/go-kit/kit/log"
)

// LogMiddleware describes a service middleware.
type LogMiddleware struct {
	Logger log.Logger
	Next   Streamer
}

// Process represents a mechanism to process streaming data.
func (lm LogMiddleware) Process(con config.Connection) (err error) {
	logger := log.With(lm.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	logger.Log(
		"Service", con.StreamingService,
		"UploadSize", con.UploadSize,
		"Topic", con.KafkaCon.Topic,
	)
	defer func(begin time.Time) {
		method := "stream.Process"
		logger.Log(
			"method", method,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = lm.Next.Process(con)
	return
}
