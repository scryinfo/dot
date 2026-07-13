CREATE TABLE users (
    id           UUID NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    first_name   VARCHAR NOT NULL DEFAULT '',
    last_name    VARCHAR NOT NULL DEFAULT '',
    display_name VARCHAR NOT NULL DEFAULT '',
    bio          TEXT DEFAULT '',
    created_at   bigint NOT NULL DEFAULT 0
);
