package service

import (
	"context"

	"gomerce/internal/model"
	"gomerce/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) Create(ctx context.Context, category *model.Category) error {
	return s.repo.Create(ctx, category)
}

func (s *CategoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *CategoryService) Update(ctx context.Context, category *model.Category) error {
	return s.repo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

