package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
)

var (
	wait, writeTimeout, readTimeout, idleTimeout time.Duration
)

func main() {
	viper.SetConfigType("ini")
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error(err.Error())
	}

	run(viper.GetInt("default.port"))
}

func run(port int) {
	var Logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	writeTimeout = time.Second * 15
	readTimeout = time.Second * 15
	idleTimeout = time.Second * 60

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		WriteTimeout: writeTimeout,
		Handler:      apiRoutes(),
	}

	Logger.Info("Ready to serve", "port", port)
	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Error("Server Initialisation", "error", err.Error())
			os.Exit(0)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	Logger.Warn("shutting down")
}
