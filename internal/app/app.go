package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"account-service/internal/config"
	"account-service/internal/handler"
	"account-service/internal/repository"
	"account-service/internal/service/auth"
	"account-service/pkg/log"
	"account-service/pkg/server"
)

// Run initializes whole application.
func Run() {
	logger := log.LoggerFromContext(context.Background())

	// Init configs
	configs, err := config.Init()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	// Init repositories
	repositories, err := repository.New(
		repository.WithPostgresStore(configs.POSTGRES.DSN))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}
	defer repositories.Close()

	authService, err := auth.New(
		auth.WithUserRepository(repositories.User),
	)
	if err != nil {
		logger.Error("ERR_INIT_AUTH_SERVICE", zap.Error(err))
		return
	}

	// Init handlers
	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			AuthService: authService,
		},
		handler.WithHTTPHandler(),
	)
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	// Init servers
	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port))

	// Run our server in a goroutine so that it doesn't block.
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}

	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		return
	}

	logger.Info("Server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // Create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("Server was successful shutdown.")
}
