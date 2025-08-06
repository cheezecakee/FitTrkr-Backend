-- +goose Up
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    playlist_id INT REFERENCES playlists(id) ON DELETE CASCADE,
    session INT,
    metadata JSONB DEFAULT '{}',
    type VARCHAR(100) NOT NULL, -- e.g., "PR_Achieved", "Workout_Completed"
    priority VARCHAR(20) NOT NULL,  -- "Legendary", "Rare", "Uncommon", "Common"
    message TEXT NOT NULL,
    pr BOOLEAN DEFAULT FALSE,  -- Store if a PR was achieved
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_logs_playlist_id ON logs(playlist_id); 
CREATE INDEX idx_logs_created_at ON logs(created_at);
CREATE INDEX idx_logs_type ON logs(type);

CREATE TRIGGER update_logs_timestamp
    BEFORE UPDATE ON logs
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE logs;                   -- Depends on sessions, users
