CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE channels_categories (
    category_id INTEGER NOT NULL,
    channel_id INTEGER NOT NULL,
    CONSTRAINT channels_categories_pks UNIQUE(category_id, channel_id)
);

ALTER TABLE channels RENAME COLUMN categories TO keywords;

CREATE OR REPLACE FUNCTION channels_search_trigger()
  RETURNS trigger AS
$$
BEGIN
  new.tsv :=
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.title, '')), 'A') ||
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.keywords, '')), 'A') ||
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.description, '')), 'D');
  RETURN new;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_categories(my_channel_id INTEGER, my_names TEXT[]) RETURNS VOID AS
$$
DECLARE
my_name TEXT;
my_category_id INTEGER;
BEGIN
    -- DELETE ALL CURRENT CATEGORIES
    DELETE FROM channels_categories WHERE channel_id=my_channel_id;
    FOREACH my_name IN ARRAY my_names
    LOOP
        BEGIN
            INSERT INTO categories(name) VALUES(LOWER(my_name)) RETURNING id INTO my_category_id;
            EXCEPTION WHEN unique_violation THEN
            SELECT id FROM categories WHERE name=LOWER(my_name) INTO my_category_id;
        END;
        BEGIN
            INSERT INTO channels_categories(channel_id, category_id) VALUES (my_channel_id, my_category_id);
            EXCEPTION WHEN unique_violation THEN
            -- do nothing
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;


DROP FUNCTION upsert_channel(text, text, text, text, text, text);

CREATE OR REPLACE FUNCTION upsert_channel(
    my_url text,
    my_title text,
    my_description text,
    my_image text,
    my_keywords text,
    my_website text)
  RETURNS integer AS
$$
DECLARE
attempts INTEGER := 1;
channel_id INTEGER;
BEGIN
    LOOP
        UPDATE channels SET 
            description=my_description, 
            title=my_title, image=my_image, 
            keywords=my_keywords, 
            website=my_website 
            WHERE url=my_url RETURNING id INTO channel_id;

        IF found THEN
            RETURN channel_id;
        END IF;

        BEGIN
            INSERT INTO channels (title, description, image, url, keywords, website) 
            VALUES (
                my_title, 
                my_description, 
                my_image, 
                my_url, 
                my_keywords, 
                my_website
            ) RETURNING id INTO channel_id;
            RETURN channel_id;
            EXCEPTION WHEN unique_violation THEN
            IF attempts = 1 THEN
               attempts := 0;
            ELSE
               RETURN 0;
            END IF;
        END;
    END LOOP;
END;
$$
LANGUAGE plpgsql;
