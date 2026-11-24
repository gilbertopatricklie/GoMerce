package repository

import (
	"context"
	"fmt"

	"gomerce/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//This struct stores database connection so that the connection could be used in other method.
type ProductRepository struct {
	db *pgxpool.Pool
}

//Creates repository object with database pool intact
func NewProductRepository (db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

//This method is to get the sum of products inside the database
func (r *ProductRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM products"

	err := r.db.QueryRow(ctx, query).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count products: %w", err)
	}

	return count, nil

}

//This method is to create product
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	query := `INSERT INTO products(
		name,
		description,
		price,
		stock
	) VALUES ($1, $2, $3, $4) RETURNING id, is_active, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.Stock).Scan(
		&product.ID, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create products: %w", err)
	}

	return nil

}

//This method is to get product by id
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `SELECT id, name, description, price, stock, is_active, created_at, updated_at
			  FROM products WHERE id = $1`

	var product model.Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}

	return &product, nil
}

//This method is to get all products
func (r *ProductRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	query := `SELECT id, name, description, price, stock, is_active, created_at, updated_at
			  FROM products ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("get all products: %w", err)
	}
	defer rows.Close()

	products := make([]model.Product, 0)

	for rows.Next() {
		var product model.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return products, nil
}

//This method update a product detail
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	query := `UPDATE products SET
			  name = $1, description = $2, price = $3, stock = $4, updated_at = NOW()
			  WHERE id = $5`

	tag, err := r.db.Exec(ctx, query, product.Name, product.Description, product.Price, product.Stock, product.ID)

	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

//This method delete a product
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}