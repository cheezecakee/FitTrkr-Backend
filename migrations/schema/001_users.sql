-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    roles TEXT[] NOT NULL DEFAULT ARRAY['user']
);

CREATE TABLE user_stats (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Body metrics
    weight NUMERIC(5,2),
    height NUMERIC(5,2),
    body_fat_percent NUMERIC(4,1),
    
    -- Workout streaks
    current_streak INT DEFAULT 0,
    longest_streak INT DEFAULT 0,
    last_workout_date DATE,
    
    -- Lifetime totals
    total_workouts INT DEFAULT 0,
    total_volume_lifted NUMERIC(10,2) DEFAULT 0,
    total_time_minutes INT DEFAULT 0,
    
    recorded_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_user_stats_timestamp
    BEFORE UPDATE ON user_stats
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE users;
