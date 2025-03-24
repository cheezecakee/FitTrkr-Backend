-- +goose Up
CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL, 
    equipment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO exercises (name, description, category, equipment) VALUES
('Push-up', 'A bodyweight exercise for chest and triceps.', 'Chest', 'Bodyweight'),
('Squat', 'A lower body exercise targeting the quadriceps and glutes.', 'Legs', 'None'),
('Deadlift', 'A full-body exercise primarily targeting the back, hamstrings, and glutes.', 'Back', 'Barbell'),
('Bench Press', 'A weightlifting exercise for the chest and triceps.', 'Chest', 'Barbell'),
('Pull-up', 'A bodyweight exercise for the back and biceps.', 'Back', 'Pull-up Bar');

-- +goose Down
DROP TABLE exercises;                      -- Standalone table, but referenced by workout_exercises
