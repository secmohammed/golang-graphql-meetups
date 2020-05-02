CREATE TABLE IF NOT EXISTS  comments(
    id BIGSERIAL PRIMARY KEY,
    body text NOT NULL,
    user_id BIGSERIAL REFERENCES users(id) ON DELETE CASCADE,
    meetup_id BIGSERIAL REFERENCES meetups(id) ON DELETE CASCADE,
    parent_id BIGSERIAL REFERENCES comments(id) ON DELETE CASCADE,
    group_id BIGSERIAL REFERENCES groups(id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
)