DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'attendance_status') THEN
        CREATE TYPE attendance_status AS ENUM
        (
            'going',
            'interested'
        );
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS  attendees(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    meetup_id BIGSERIAL REFERENCES meetups (id) ON DELETE CASCADE NOT NULL,
    status attendance_status,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL

);