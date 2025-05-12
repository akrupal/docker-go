package handler

import (
	"api-server/internal/database"
	"api-server/internal/restapi/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type handler struct {
	service service.Service
}

type Handler interface {
	AddProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	GetAllProducts(w http.ResponseWriter, r *http.Request)
}

func NewHandler(service service.Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.WithContext(ctx)
	var product database.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		logger.Error("Failed to decode user input")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Failed to decode user input")
	}
	if product.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Product name missing")
	}
	pk := h.service.AddProduct(product)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf("Product added with product key %v", pk))
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	user := mux.Vars(r)
	id, err := strconv.Atoi(user["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("wrong id passed")
	}

	product := h.service.GetProduct(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products := h.service.GetAllProducts()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
