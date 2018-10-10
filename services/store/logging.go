package store

import (
	"octopus/config"
	"time"

	log "github.com/go-kit/kit/log"
)

// LogMiddleware describes a service middleware.
type LogMiddleware struct {
	Logger log.Logger
	Next   Store
}

// Minio stores files in Minio S3 compatable storage.
func (lm LogMiddleware) Minio(config config.Minio, b []byte, topic string) (object string, err error) {
	defer func(begin time.Time) {
		method := "store.Minio"
		logger := log.With(lm.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
		_ = logger.Log(
			"method", method,
			"topic", topic,
			"object", object,
			"size", len(b),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	object, err = lm.Next.Minio(config, b, topic)
	return
}
