package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/mntwo/tasklab/internal/application"
	"github.com/mntwo/tasklab/internal/config"
	"github.com/mntwo/tasklab/internal/log"
	"go.uber.org/zap"
)

func Run(apps ...application.Application) {
	ctx := context.Background()
	log.Info(ctx, "tasklab running", zap.Any("application", config.GetApplication()))

	wg := sync.WaitGroup{}
	for _, app := range apps {
		a := app
		wg.Add(1)
		go func(app application.Application) {
			defer wg.Done()
			log.Info(ctx, "starting application", zap.String("name", app.GetName()))
			if err := app.Start(); err != nil {
				if errors.Is(err, http.ErrServerClosed) || errors.Is(err, application.ErrApplicationClosed) {
					log.Info(ctx, "application closed", zap.String("name", app.GetName()))
					return
				}
				log.Fatal(ctx, "failed to start application", zap.Error(err), zap.String("name", app.GetName()))
			}
		}(a)
	}

	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	sig := <-signals
	switch sig {
	case syscall.SIGTERM:
		log.Fatal(ctx, fmt.Sprintf("force exit received signal=%s", sig))
	case syscall.SIGHUP, syscall.SIGINT:
		log.Info(ctx, fmt.Sprintf("graceful shutdown received signal=%s\n", sig))
		for _, app := range reverse(apps) {
			if err := app.Stop(); err != nil {
				log.Fatal(ctx, "failed to stop application", zap.Error(err), zap.String("name", app.GetName()))
			}
		}
		wg.Wait()
		log.Info(ctx, "tasklab stopped")
	}
}

func reverse(apps []application.Application) []application.Application {
	for i := 0; i < len(apps)/2; i++ {
		j := len(apps) - i - 1
		apps[i], apps[j] = apps[j], apps[i]
	}
	return apps
}
