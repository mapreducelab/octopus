package process

import (
	"time"

	log "github.com/go-kit/kit/log"
)

// LogMiddleware describes a service middleware.
type LogMiddleware struct {
	Logger log.Logger
	Next   Processor
}

// Process messages from streaming services.
func (lm LogMiddleware) Process(b []byte) (oBytes []byte, err error) {
	defer func(begin time.Time) {
		method := "store.Process"
		logger := log.With(lm.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
		_ = logger.Log(
			"method", method,
			"size", len(b),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	oBytes, err = lm.Next.Process(b)
	return
}
