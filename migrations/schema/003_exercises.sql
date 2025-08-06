-- +goose Up

-- Exercise Categories
CREATE TABLE exercise_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Equipment
CREATE TABLE equipment (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Exercise Types (like tags)
CREATE TABLE training_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

-- Exercises
CREATE TABLE exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    category_id INT NOT NULL REFERENCES exercise_categories(id),
    equipment_id INT NOT NULL REFERENCES equipment(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Trigger to auto-update updated_at
CREATE TRIGGER update_exercises_timestamp
    BEFORE UPDATE ON exercises
    FOR EACH ROW EXECUTE FUNCTION update_timestamp();

-- Muscle Groups
CREATE TABLE muscle_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Join Table: Exercise ↔ Muscle Groups
CREATE TABLE exercise_muscles (
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    muscle_group_id INT NOT NULL REFERENCES muscle_groups(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, muscle_group_id)
);

-- Join Table: Exercise ↔ Types (many-to-many)
CREATE TABLE exercise_training_types (
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    type_id INT NOT NULL REFERENCES training_types(id) ON DELETE CASCADE,
    PRIMARY KEY (exercise_id, type_id)
);

-- Indexes for fast lookups
CREATE INDEX idx_exercise_muscles_exercise_id ON exercise_muscles(exercise_id);
CREATE INDEX idx_exercise_muscles_muscle_group_id ON exercise_muscles(muscle_group_id);
CREATE INDEX idx_exercise_training_types_exercise_id ON exercise_training_types(exercise_id);
CREATE INDEX idx_exercise_training_types_type_id ON exercise_training_types(type_id);

-- +goose Down
DROP TABLE IF EXISTS exercise_types_junction;
DROP TABLE IF EXISTS exercise_muscles;
DROP TABLE IF EXISTS exercises;
DROP TABLE IF EXISTS muscle_groups;
DROP TABLE IF EXISTS exercise_categories;
DROP TABLE IF EXISTS equipment;
DROP TABLE IF EXISTS exercise_types;

-- TODO: Consider adding equipment_muscles table
-- To associate equipment with muscle groups
-- Useful for recommendation, filtering, or workout generation
