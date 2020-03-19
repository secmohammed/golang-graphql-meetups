CREATE TABLE category_meetup(
    meetup_id BIGSERIAL REFERENCES meetups (id) ON DELETE CASCADE NOT NULL,
    category_id BIGSERIAL REFERENCES categories (id) ON DELETE CASCADE NOT NULL
)