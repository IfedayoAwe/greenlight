CREATE TABLE IF NOT EXISTS users_profile (
    user_id bigint UNIQUE NOT NULL REFERENCES users ON DELETE CASCADE,
    image_path text NOT NULL
);