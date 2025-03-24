-- +goose Up
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,
    plan_id INT,  -- Nullable, as a session may not follow a specific plan
    name VARCHAR(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,  -- Nullable until workout is completed
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE SET NULL
);

CREATE TABLE sessions_exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id INT NOT NULL,
    exercise_id INT NOT NULL,
    "order" INT NOT NULL,  -- The sequence of exercises within the session
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    UNIQUE (session_id, exercise_id)  -- Each exercise appears once per session
);

CREATE TABLE exercise_sets (
    id SERIAL PRIMARY KEY,
    session_exercise_id UUID NOT NULL REFERENCES sessions_exercises(id) ON DELETE CASCADE,
    set_number INT NOT NULL,
    reps INT,
    weight NUMERIC(6,2),    -- Allows for weights like 185.25 lbs
    duration INT,           -- Duration in seconds
    distance NUMERIC(6,2),  -- For exercises like running, swimming, etc.
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_exercise_id) REFERENCES sessions_exercises(id) ON DELETE CASCADE,
    UNIQUE (session_exercise_id, set_number)  -- Sets should be uniquely numbered within an exercise
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_exercises_session_id ON sessions_exercises(session_id);

-- +goose Down
DROP TABLE exercise_sets;           -- Depends on sessions_exercises
DROP TABLE sessions_exercises;      -- Depends on sessions, plans_exercises
DROP TABLE sessions;                -- Depends on plans, users
