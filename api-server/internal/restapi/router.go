package restapi

import (
	"api-server/internal/configurator"
	"api-server/internal/database"
	"api-server/internal/restapi/handler"
	"api-server/internal/restapi/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Server    *http.Server
	Waitgroup sync.WaitGroup
}

func Start(config *configurator.SqlConfig) *Server {

	db, err := database.CreateNewResourceRepository(context.Background(), config)
	if err != nil {
		log.Fatal("failed to create a new resource repository")
	}
	newService := service.NewService(db)
	newHandler := handler.NewHandler(newService)

	router := mux.NewRouter()

	router.HandleFunc("/add-product", newHandler.AddProduct).Methods(http.MethodPost)

	router.HandleFunc("/get-product/{id}", newHandler.GetProduct).Methods(http.MethodGet)

	router.HandleFunc("/get-all-products", newHandler.GetAllProducts).Methods(http.MethodGet)

	Rhandler := handlers.RecoveryHandler()(router)
	server := &Server{
		Server: &http.Server{
			Handler: Rhandler,
			Addr:    fmt.Sprintf(":%d", 7777),
		},
	}
	server.Waitgroup.Add(1)
	go func() {
		logrus.Debug("API server started listening")
		err := server.Server.ListenAndServe()
		if err != nil {
			logrus.Debug("API server unexpected termination")
		}
		server.Waitgroup.Done()
	}()

	return server
}

func (Server *Server) Stop() {
	logrus.Info("Stopping API server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := Server.Server.Shutdown(ctx); err != nil {
			logrus.Fatal("Failed to shutdown API server")
		}
	}()

	Server.Waitgroup.Wait()
	logrus.Info("API server shutting down")
}
