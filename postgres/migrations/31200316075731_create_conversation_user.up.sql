CREATE TABLE conversation_user(
    conversation_id BIGSERIAL REFERENCES conversations (id) ON DELETE CASCADE NOT NULL,
    user_id BIGSERIAL REFERENCES users (id) ON DELETE CASCADE NOT NULL
)