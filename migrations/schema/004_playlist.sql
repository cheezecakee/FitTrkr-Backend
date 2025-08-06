-- +goose Up
CREATE TABLE playlists (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT FALSE,
    last_worked BOOLEAN DEFAULT FALSE,
    visibility VARCHAR(20) NOT NULL DEFAULT 'private', -- 'private', 'public', 'unlisted'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_user_playlist_title UNIQUE (user_id, title)
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE playlist_tags (
    playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (playlist_id, tag_id)
);

CREATE TABLE exercise_blocks (
    id SERIAL PRIMARY KEY,
    playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    name VARCHAR(100),
    block_type VARCHAR(20) DEFAULT 'regular', -- 'regular', 'superset', 'circuit', 'dropset'
    block_order INT NOT NULL,
    rest_after_block_seconds INT DEFAULT 60 -- Rest after completing entire block
);

CREATE TABLE exercise_configs (
    id SERIAL PRIMARY KEY,
    
    -- Strength fields
    sets INT,
    reps_min INT,
    reps_max INT,
    weight NUMERIC(6,2),
    
    -- Cardio fields  
    duration_seconds INT,
    distance NUMERIC(6,2),
    target_pace NUMERIC(5,2), -- minutes per mile/km
    target_heart_rate INT,
    incline NUMERIC(4,1),
    
    -- Rest is handled at exercise level (not block level)
    rest_seconds INT DEFAULT 60,
    tempo INT[] CHECK (array_length(tempo, 1) = 4), -- [eccentric, pause, concentric, pause]
    
    -- Common fields
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE playlist_exercises (
    id SERIAL PRIMARY KEY,
    playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    block_id INT NOT NULL REFERENCES exercise_blocks(id) ON DELETE CASCADE, -- Made NOT NULL
    config_id INT NOT NULL REFERENCES exercise_configs(id) ON DELETE CASCADE,
    exercise_order INT NOT NULL DEFAULT 0, -- Order within the block
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_playlist_exercises_playlist_id ON playlist_exercises(playlist_id);
CREATE INDEX idx_playlist_exercises_exercise_id ON playlist_exercises(exercise_id);
CREATE INDEX idx_playlist_exercises_block_id ON playlist_exercises(block_id);
CREATE INDEX idx_playlist_exercises_config_id ON playlist_exercises(config_id);
CREATE INDEX idx_exercise_blocks_playlist_id ON exercise_blocks(playlist_id);

-- Triggers
CREATE TRIGGER update_playlists_timestamp
    BEFORE UPDATE ON playlists
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_playlist_exercises_timestamp
    BEFORE UPDATE ON playlist_exercises
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER update_exercise_configs_timestamp
    BEFORE UPDATE ON exercise_configs
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- +goose Down
DROP TABLE IF EXISTS playlist_exercises;
DROP TABLE IF EXISTS playlist_tags;
DROP TABLE IF EXISTS exercise_configs;
DROP TABLE IF EXISTS exercise_blocks;
DROP TABLE IF EXISTS playlists;
DROP TABLE IF EXISTS tags;
