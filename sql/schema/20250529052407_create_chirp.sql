-- +goose Up
-- +goose StatementBegin
CREATE TABLE chirps (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT Now(),
    updated_at TIMESTAMP NOT NULL Default Now(),
    body TEXT UNIQUE,
    user_id uuid REFERENCES users ON DELETE CASCADE,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chirps;
-- +goose StatementEnd
