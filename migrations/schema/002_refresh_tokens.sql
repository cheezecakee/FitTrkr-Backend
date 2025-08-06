-- +goose Up
CREATE TABLE refresh_tokens (
    token VARCHAR(64) PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL DEFAULT now(), 
    updated_at TIMESTAMP NOT NULL DEFAULT now(), 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    is_revoked BOOL NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL, 
    revoked_at TIMESTAMP
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE TRIGGER update_refresh_tokens_timestamp
    BEFORE UPDATE ON refresh_tokens
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE refresh_tokens;
