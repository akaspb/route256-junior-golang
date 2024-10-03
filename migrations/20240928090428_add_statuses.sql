-- +goose Up
CREATE TABLE statuses (
    order_id BIGINT PRIMARY KEY,
    "value" VARCHAR(255) NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders ON DELETE CASCADE
);

-- +goose Down
DROP TABLE statuses;
