ALTER TABLE macro_entries ADD COLUMN serving_weight NUMERIC(7,1) NOT NULL DEFAULT 0;
ALTER TABLE macro_entries ADD COLUMN serving_unit VARCHAR(8) NOT NULL DEFAULT 'g';
