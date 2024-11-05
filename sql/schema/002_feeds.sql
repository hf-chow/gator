-- +goose Up
CREATE TABLE feeds (
    id              uuid PRIMARY KEY,
    created_at      timestamp NOT NULL,
    updated_at      timestamp NOT NULL,
    name            varchar(40) NOT NULL,
    url             varchar(40) NOT NULL,
    user_id         uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) on DELETE CASCADE
);
-- +goose Down
DROP TABLE feeds;
