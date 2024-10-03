-- +goose Up
CREATE TABLE packs (
    order_id BIGINT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL CONSTRAINT non_empty_name CHECK (LENGTH(name) > 0),
    cost NUMERIC(12,2) NOT NULL,
    max_order_weight NUMERIC(7, 3),
    FOREIGN KEY (order_id) REFERENCES orders ON DELETE CASCADE
);

-- +goose Down
DROP TABLE packs;
