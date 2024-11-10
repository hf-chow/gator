-- +goose Up
CREATE TABLE feed_follows (
    id              uuid PRIMARY KEY,
    created_at      timestamp NOT NULL,
    updated_at      timestamp NOT NULL,
    user_id         uuid NOT NULL,
    feed_id         uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) on DELETE CASCADE,
    UNIQUE(user_id, feed_id)
);
-- +goose Down
DROP TABLE feed_follows;