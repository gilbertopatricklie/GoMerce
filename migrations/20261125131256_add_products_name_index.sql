-- +goose Up
CREATE INDEX idx_products_name ON products (name);

-- +goose Down
DROP INDEX idx_products_name;
