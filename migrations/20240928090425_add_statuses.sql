-- +goose Up
CREATE TABLE statuses (
    id BIGSERIAL PRIMARY KEY,
    "value" VARCHAR(255) NOT NULL,
    "time" TIMESTAMP WITH TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS statuses;
