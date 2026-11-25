-- +goose Up
CREATE TABLE categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE products
ADD COLUMN category_id BIGINT;

ALTER TABLE products
ADD CONSTRAINT fk_products_category
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE SET NULL;

CREATE INDEX idx_products_category_id
ON products(category_id);

-- +goose Down
ALTER TABLE products
DROP COLUMN category_id;

DROP TABLE categories;
