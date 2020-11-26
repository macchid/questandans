package utils

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func LogStart(parent log.Logger, method string) (log.Logger, func(time.Time)) {
	logger := log.With(parent, "method", method)
	level.Info(logger).Log("msg", "Starting method")

	return logger, func(start time.Time) {
		level.Info(logger).Log("msg", "Method completed.", "took", time.Since(start))
	}
}
