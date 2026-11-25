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

//Passes method result to repository / product_repository.go
func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.GetAll(ctx)
}
//Passes method result to repository / product_repository.go
func (s *ProductService) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}
//Passes method result to repository / product_repository.go
func (s *ProductService) Update(ctx context.Context, product *model.Product) error {
	return s.repo.Update(ctx, product)
}
//Passes method result to repository / product_repository.go
func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
