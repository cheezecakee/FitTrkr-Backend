-- +goose Up
CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE plans_exercises (
    id SERIAL PRIMARY KEY,
    plan_id INT NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sets INT NOT NULL,
    reps INT[] NOT NULL CHECK (array_length(reps, 1) = 2),  -- Ensures exactly 2 values
    weight FLOAT NOT NULL,
    interval INT NOT NULL DEFAULT 0, -- Store as seconds
    rest INT NOT NULL DEFAULT 60, -- Store as seconds
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE plans;                       -- Depends on users
DROP TABLE plans_exercises;             -- Depends on plans, exercises
