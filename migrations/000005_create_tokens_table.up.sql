CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    user_ip text NOT NULL,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);