ALTER TABLE programs
    ADD COLUMN workouts_per_week INT NOT NULL DEFAULT 4 CHECK (workouts_per_week > 0);
