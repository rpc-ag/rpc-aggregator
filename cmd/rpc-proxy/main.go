package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rpc-ag/rpc-proxy/internal/config"
	"github.com/rpc-ag/rpc-proxy/internal/webserver"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	var configFilePath string
	flag.StringVar(&configFilePath, "config", "./config.yaml", "config file to load")
	flag.Parse()
	conf, err := config.Read(configFilePath)
	if err != nil {
		logger.Panic("failed to load config", zap.Error(err))
	}
	logger.Info("config loaded", zap.Any("conf", conf))

	server, err := webserver.New(conf, logger)
	if err != nil {
		logger.Panic("failed to start webserver", zap.Error(err))
	}

	logger.Info("Starting RPC proxy...")

	go func() {
		er := server.Run()
		if er != nil {
			panic(er)
		}
	}()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)
	<-interrupt

	// Gracefully shutdown server
	logger.Info("Shutting down RPC proxy...")
	if err := server.Close(); err != nil {
		logger.Error("Error shutting down server")
	}
}
