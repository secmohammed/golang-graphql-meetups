CREATE TABLE IF NOT EXISTS category_group (
    group_id BIGSERIAL REFERENCES groups (id) ON DELETE CASCADE NOT NULL,
    category_id BIGSERIAL REFERENCES categories (id) ON DELETE CASCADE NOT NULL
)