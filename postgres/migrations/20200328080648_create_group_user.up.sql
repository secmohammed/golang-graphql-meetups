DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'member_status') THEN
        CREATE TYPE member_status AS ENUM
        (
            'admin',
            'moderator',
            'member'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS  group_user(
    group_id BIGSERIAL REFERENCES groups (id) ON DELETE CASCADE NOT NULL,
    user_id BIGSERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    type member_status
)