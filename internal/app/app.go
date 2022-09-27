package app

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	gw "github.com/Teamlix/proto/gen/go/user_service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/teamlix/api-gateway/internal/pkg/config"
	log "github.com/teamlix/api-gateway/internal/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(configPath string) error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return err
	}

	logger, err := log.NewLogger()
	if err != nil {
		return nil
	}

	httpAddr := fmt.Sprintf("%s:%s", cfg.Http.Host, cfg.Http.Port)
	go func() {
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		err := gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.Grpc.Client.User, opts)
		if err != nil {
			errCh <- err
		}

		err = http.ListenAndServe(httpAddr, mux)
		if err != nil {
			errCh <- err
		}
	}()

	logger.Infof("HTTP server strarted on: %s", httpAddr)

	// listen to os signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Infof("OS signal: %s", sig.String())

	case err := <-errCh:
		logger.Errorf("Got error: %v", err)
	}

	logger.Info("http listener closed")

	return nil
}
