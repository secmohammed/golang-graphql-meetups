CREATE TABLE IF NOT EXISTS  roles(
    id BIGSERIAL PRIMARY KEY,
    permissions text NOT NULL,
    user_id BIGSERIAL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);
CREATE TABLE IF NOT EXISTS role_user(
    user_id BIGSERIAL REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    role_id BIGSERIAL REFERENCES roles(id) ON DELETE CASCADE NOT NULL
);