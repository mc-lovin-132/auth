package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mc-lovin-132/auth/config"
	"github.com/mc-lovin-132/auth/internal/app"

	"go.uber.org/zap"
)

func main() {
	// cfg
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	app := app.New(cfg, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		cancel()
	}()

	err = app.Start(ctx)
	if err != nil {
		logger.Fatal("error running server", zap.Error(err))
	}
}
