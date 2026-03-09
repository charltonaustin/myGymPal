ALTER TABLE programs
    DROP CONSTRAINT programs_user_id_fkey,
    ADD CONSTRAINT programs_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES users(id);
