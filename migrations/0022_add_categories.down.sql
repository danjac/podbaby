DROP FUNCTION add_categories (INTEGER, TEXT[]);
DROP TABLE channels_categories;
DROP TABLE categories;

ALTER TABLE channels RENAME COLUMN keywords TO categories;

CREATE OR REPLACE FUNCTION channels_search_trigger()
  RETURNS trigger AS
$$
BEGIN
  new.tsv :=
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.title, '')), 'A') ||
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.categories, '')), 'A') ||
    SETWEIGHT(TO_TSVECTOR(COALESCE(new.description, '')), 'D');
  RETURN new;
END;
$$ LANGUAGE plpgsql;


DROP FUNCTION upsert_channel(text, text, text, text, text, text);

CREATE OR REPLACE FUNCTION upsert_channel(
    my_url text,
    my_title text,
    my_description text,
    my_image text,
    my_categories text,
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
            categories=my_categories, 
            website=my_website 
            WHERE url=my_url RETURNING id INTO channel_id;

        IF found THEN
            RETURN channel_id;
        END IF;

        BEGIN
            INSERT INTO channels (title, description, image, url, categories, website) 
            VALUES (
                my_title, 
                my_description, 
                my_image, 
                my_url, 
                my_categories, 
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
