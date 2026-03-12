ALTER TABLE template_exercises ADD COLUMN block VARCHAR(20) NOT NULL DEFAULT 'main';
ALTER TABLE session_exercises ADD COLUMN block VARCHAR(20) NOT NULL DEFAULT 'main';
