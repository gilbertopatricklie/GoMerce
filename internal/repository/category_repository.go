package repository

import (
	"context"
	"fmt"

	"gomerce/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

//Method to create a category
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, category.Name).Scan(&category.ID, &category.CreatedAt)

	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}

	return nil
}

//Method to get ALL categories
func (r *CategoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("get all categories: %w", err)
	}

	defer rows.Close()

	categories := make([]model.Category, 0)

	for rows.Next() {
		var category model.Category

		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}

	return categories, nil
}