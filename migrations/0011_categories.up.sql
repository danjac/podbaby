CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE channel_categories (
    channel_id INTEGER NOT NULL REFERENCES channels(id),
    category_id INTEGER NOT NULL REFERENCES categories(id)
);

ALTER TABLE channel_categories ADD CONSTRAINT unique_channel_categories UNIQUE(channel_id, category_id);

CREATE UNIQUE INDEX ON categories (name ASC);

CREATE INDEX ON channel_categories (channel_id ASC);
CREATE INDEX ON channel_categories (category_id ASC);
