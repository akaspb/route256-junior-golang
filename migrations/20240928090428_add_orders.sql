-- +goose Up
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL,
    status_id BIGINT NOT NULL REFERENCES statuses ON DELETE CASCADE UNIQUE,
    weight NUMERIC(7, 3) NOT NULL CONSTRAINT non_negative_weight CHECK (NOT (weight < 0)),
    cost MONEY NOT NULL CONSTRAINT non_negative_cost CHECK (NOT (cost < 0::MONEY))
);

-- +goose Down
DROP TABLE IF EXISTS orders;
