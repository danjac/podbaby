ALTER TABLE categories DROP CONSTRAINT categories_parent_id_fk;
ALTER TABLE categories DROP COLUMN parent_id;

DROP FUNCTION add_categories(INTEGER, TEXT[]);

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


