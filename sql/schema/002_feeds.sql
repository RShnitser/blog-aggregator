-- +goose Up
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT UNIQUE NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL DEFAULT gen_random_uuid(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;