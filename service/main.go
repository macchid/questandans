package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	qsrv "github.com/macchid/questandans/questions/implementation"
	repo "github.com/macchid/questandans/repository"
	"github.com/macchid/questandans/transport"
	httpt "github.com/macchid/questandans/transport/http"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
		mongoURI = flag.String("mongo.uri", "localhost:27017", "Mongo DB liste address")
	)

	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "question",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "Service started")
	defer level.Info(logger).Log("msg", "Service ended")

	timeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := repo.NewRepository(timeout, logger, *mongoURI)
	if err != nil {
		// TODO: loguear el error.
		os.Exit(1)
	}

	qs := qsrv.NewService(db, logger)

	endpoints := transport.MakeEndpoints(qs)

	r := httpt.NewService(mux.NewRouter(), endpoints, logger)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("Received the signal %s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: r,
		}
		errs <- server.ListenAndServe()
	}()

	level.Error(logger).Log("exit", <-errs)
}
