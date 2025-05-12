package main

import (
	"api-server/internal/configurator"
	"api-server/internal/restapi"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	maxOpenConn, err := strconv.Atoi(os.Getenv("MaxOpenConns"))
	if err != nil {
		log.Fatal("Invalid maxOpenConn")
	}

	maxIdleConn, err := strconv.Atoi(os.Getenv("MaxIdleConns"))
	if err != nil {
		log.Fatal("Invalid maxIdleConn")
	}

	connMaxLifetime, err := strconv.Atoi(os.Getenv("ConnMaxLifetime"))
	if err != nil {
		log.Fatal("invalid connMaxLifetime")
	}

	config := &configurator.SqlConfig{
		DbDriver:        os.Getenv("DbDriver"),
		DbHost:          os.Getenv("DbHost"),
		DbPort:          os.Getenv("DbPort"),
		DbName:          os.Getenv("DbName"),
		DbUsername:      os.Getenv("DbUsername"),
		DbPassword:      os.Getenv("DbPassword"),
		DbSslMode:       os.Getenv("DbSslMode"),
		MaxOpenConns:    maxOpenConn,
		MaxIdleConns:    maxIdleConn,
		ConnMaxLifetime: int64(connMaxLifetime),
	}

	defer restapi.Start(config).Stop()

	sigChan := make(chan os.Signal, 1)
	// Catching SIGINT (Ctrl+C) AND SIGTERM (Ctrl+/) for graceful shutdown. SIGKILL or SIGQUIT will not be caught
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	logrus.WithField("signal", <-sigChan).Info("Received shutdown signal, shutting down...")
	logrus.Info("Resource distribution API server terminated without error")
}
