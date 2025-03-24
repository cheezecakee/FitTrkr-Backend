-- +goose Up
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL DEFAULT now(), 
    updated_at TIMESTAMP NOT NULL DEFAULT now(), 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    is_revoked BOOL NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL, 
    revoked_at TIMESTAMP
);

-- +goose Down
DROP TABLE refresh_tokens;                 -- Depends on users
