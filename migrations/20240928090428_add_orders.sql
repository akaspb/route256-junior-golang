-- +goose Up
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    status_id BIGINT NOT NULL REFERENCES statuses ON DELETE CASCADE UNIQUE,
    weight NUMERIC(7, 3) NOT NULL,
    cost NUMERIC(12,2) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;
