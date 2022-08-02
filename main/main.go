package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus" // MIT

	"github.com/dontang97/AU0/router"
	"github.com/dontang97/AU0/util/auth"
	"github.com/dontang97/AU0/util/cryptor"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	// dbhost := flag.String("db-host", "db", "the database host")
	// dbport := flag.Int("db-port", 5432, "the database port")

	addr := flag.String("address", ":9998", "listening address")
	logLevel := flag.String("log-level", "DEBUG", "log level")
	flag.Parse()

	log := logrus.New()
	initLogger(log, *logLevel, os.Stdout)

	/*
		_ui := ui.New()
		_ui.Connect(*dbhost, *dbport)
		defer _ui.Disconnect()
	*/

	crypt, err := cryptor.NewHTCTokenCryptor(
		[]byte("LIjhBWwR1A9BjiBCdsMs0KiZ3x50Ce9auGFBqqj69Q4="),
		[]byte("tKIlYqe00NtAuXhDy1UfYQ=="),
	)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	auther, err := auth.New(
		"dev-don.us.auth0.com",
		"ca8S2VKCzgdVFadAswRqmymVmQNIuzTN",
		"F9XVlojlGKSsW-CSQ7jJ8OCltXk6A0FPCNGLIitafzHdGgXlp-dbPCR5qYfCcXsV",
	)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:         *addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router.Route(log, nil, crypt, auther),
	}
	go func() {
		log.Debug("Server starts...")
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// block
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err)
	}

	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func initLogger(log *logrus.Logger, level string, out io.Writer) {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(out)

	switch level {
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "WARN":
		log.SetLevel(logrus.WarnLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.DebugLevel)
	}
}
