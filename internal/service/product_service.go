package service

import (
	"context"

	"gomerce/internal/model"
	"gomerce/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) error {
	return s.repo.Create(ctx, product)
}