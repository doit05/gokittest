package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"my/profilesvc"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
		logpath  = flag.String("logpath", "log.log", "log path")
	)
	flag.Parse()

	logfile, err := os.OpenFile(*logpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("creating %s failed: %s", *logpath, err))
	}
	logfile.Sync()
	defer logfile.Close()
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(logfile)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s profilesvc.Service
	{
		s = profilesvc.NewInmemService()
		s = profilesvc.LoggingMiddleware(logger)(s)
	}

	var h http.Handler
	{
		h = profilesvc.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
