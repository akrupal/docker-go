package service

import "api-server/internal/database"

type service struct {
	repository database.DatabaseInterface
}

type Service interface {
	AddProduct(product database.Product) int
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
