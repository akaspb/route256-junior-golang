-- +goose Up
CREATE TABLE packaging (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL CONSTRAINT non_empty_name CHECK (LENGTH(name) > 0),
    cost MONEY NOT NULL CONSTRAINT non_negative_cost CHECK (NOT (cost < 0::MONEY)),
    max_order_weight NUMERIC(7, 3) CONSTRAINT positive_max_order_weight CHECK (max_order_weight > 0)
);

-- +goose Down
DROP TABLE IF EXISTS packaging;
