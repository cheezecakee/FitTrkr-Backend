-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password_hash TEXT NOT NULL DEFAULT 'unset',
    email TEXT UNIQUE NOT NULL,
    age INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_premium BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL DEFAULT now(), 
    updated_at TIMESTAMP NOT NULL DEFAULT now(), 
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, 
    is_revoked BOOL NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL, 
    revoked_at TIMESTAMP
);

CREATE TABLE workouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO exercises (name) VALUES
('Push-up'),
('Squat'),
('Deadlift'),
('Bench Press'),
('Pull-up');

CREATE TABLE workout_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id UUID NOT NULL,
    exercise_id UUID NOT NULL,
    sets INT NOT NULL,
    reps_min INT NOT NULL,
    reps_max INT NOT NULL,
    weight FLOAT NOT NULL,
    interval INT NOT NULL DEFAULT 0, -- Store as seconds
    rest INT NOT NULL DEFAULT 60, -- Store as seconds
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);

CREATE TABLE workout_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    ended_at TIMESTAMP,  -- Null until the session ends
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT one_active_session_per_user
        EXCLUDE USING btree (user_id WITH =, is_active WITH =)
);

CREATE INDEX idx_btree_workout_session ON workout_session (user_id, is_active);

CREATE TABLE exercise_session (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_session_id UUID NOT NULL REFERENCES workout_session(id) ON DELETE CASCADE,
    workout_exercise_id UUID NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    ended_at TIMESTAMP,  -- Null until finished or skipped
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT one_active_session_per_workout
        EXCLUDE USING btree (workout_session_id WITH =, is_active WITH =)
);

CREATE INDEX idx_btree_exercise_session ON exercise_session (workout_session_id, is_active);

CREATE TABLE workout_exercise_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_session_id UUID NOT NULL REFERENCES exercise_session(id) ON DELETE CASCADE,
    set_number INT NOT NULL,  -- Tracks which set this log belongs to
    reps INT NOT NULL DEFAULT 0,
    weight DECIMAL(5,2) NOT NULL DEFAULT 0.0,
    interval_sec INT NOT NULL DEFAULT 0,  -- Rest time between sets
    notes TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE session_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_session_id UUID NOT NULL REFERENCES workout_session(id) ON DELETE CASCADE,
    workout_exercise_logs_id UUID REFERENCES workout_exercise_logs(id) ON DELETE CASCADE,
    log_type TEXT NOT NULL,  -- e.g., "PR_Achieved", "Workout_Completed"
    log_priority TEXT NOT NULL,  -- "Legendary", "Rare", "Uncommon", "Common"
    log_message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE session_logs;                   -- Depends on workout_session, workout_exercise_logs, users
DROP TABLE workout_exercise_logs;          -- Depends on workout_exercise_session
DROP TABLE exercise_session;               -- Depends on workout_session, workout_exercises
DROP TABLE workout_session;                -- Depends on workouts, users
DROP TABLE workout_exercises;              -- Depends on workouts, exercises
DROP TABLE exercises;                      -- Standalone table, but referenced by workout_exercises
DROP TABLE workouts;                       -- Depends on users
DROP TABLE refresh_tokens;                 -- Depends on users
DROP TABLE users;                          -- No dependencies, can be dropped last
