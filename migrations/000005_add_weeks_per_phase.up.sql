ALTER TABLE programs
    ADD COLUMN weeks_per_phase INT NOT NULL DEFAULT 8 CHECK (weeks_per_phase > 0);
