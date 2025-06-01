-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT Now(),
    updated_at TIMESTAMP NOT NULL Default Now(),
    email TEXT UNIQUE,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
