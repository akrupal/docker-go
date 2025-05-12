package service

import "api-server/internal/database"

type service struct {
	repository database.DatabaseInterface
}

type Service interface {
	AddProduct(product database.Product) int
	GetProduct(id int) database.Product
	GetAllProducts() []database.Product
}

func NewService(repository database.DatabaseInterface) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) AddProduct(product database.Product) int {
	pk := s.repository.AddProduct(product)
	return pk
}

func (s *service) GetProduct(id int) database.Product {
	product := s.repository.GetProduct(id)
	return product
}

func (s *service) GetAllProducts() []database.Product {
	products := s.repository.GetAllProducts()
	return products
}
