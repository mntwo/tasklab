package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mntwo/tasklab/internal/config"
	"github.com/mntwo/tasklab/internal/log"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	log.Info(ctx, "starting tasklab",
		zap.Any("application", config.GetApplication()),
		zap.Any("http_server", config.GetHttpServer()),
		zap.Any("log", config.GetLog()),
		zap.Any("mysql", config.GetMySQL()),
		zap.Any("redis", config.GetRedis()),
	)
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	sig := <-signals
	switch sig {
	case syscall.SIGTERM:
		// force exit
		log.Fatal(ctx, fmt.Sprintf("force exit received signal=%s", sig))
	case syscall.SIGHUP, syscall.SIGINT:
		// graceful shutdown
		log.Info(ctx, fmt.Sprintf("graceful shutdown received signal=%s\n", sig))
	}
}
