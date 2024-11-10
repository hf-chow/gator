-- +goose Up
CREATE TABLE posts (
    id              uuid PRIMARY KEY,
    created_at      timestamp NOT NULL,
    updated_at      timestamp NOT NULL,
    title           varchar(40) NOT NULL,
    url             varchar(40) UNIQUE NOT NULL,
    description     varchar(40) NOT NULL,
    published_at    timestamp NOT NULL,
    feed_id         UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE posts;
