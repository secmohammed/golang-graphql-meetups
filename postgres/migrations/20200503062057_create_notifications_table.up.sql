DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notifiable_type') THEN
        CREATE TYPE notifiable_type AS ENUM
        (
            'reply_created',
            'comment_created',
            'meetup_created',
            'meetup_reminder',
            'meetup_shared_to_group'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS  notifications(
    id BIGSERIAL PRIMARY KEY,
    notifiable_type notifiable_type,
    notifiable_id BIGSERIAL,
    user_id BIGSERIAL REFERENCES users(id) ON DELETE CASCADE,
    read_at timestamp with time zone,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
)