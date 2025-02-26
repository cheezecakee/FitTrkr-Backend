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

-- +goose Down
DROP TABLE refresh_tokens;
DROP TABLE workout_exercises;
DROP TABLE exercises;
DROP TABLE workouts;
DROP TABLE users;
