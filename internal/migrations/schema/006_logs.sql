-- +goose Up
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INT NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    metadata JSONB DEFAULT '{}',
    type VARCHAR(100) NOT NULL, -- e.g., "PR_Achieved", "Workout_Completed"
    priority VARCHAR(20) NOT NULL,  -- "Legendary", "Rare", "Uncommon", "Common"
    message TEXT NOT NULL,
    pr BOOLEAN DEFAULT FALSE,  -- Store if a PR was achieved
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_logs_plan_id ON logs(plan_id);
CREATE INDEX idx_logs_created_at ON logs(created_at);
CREATE INDEX idx_logs_type ON logs(type);

-- +goose Down
DROP TABLE logs;                   -- Depends on sessions, users
