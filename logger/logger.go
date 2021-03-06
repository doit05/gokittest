package logger

import (
	"github.com/go-kit/kit/log"
	"os"
)

var Logger log.Logger

func init() {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	Logger = logger
}
