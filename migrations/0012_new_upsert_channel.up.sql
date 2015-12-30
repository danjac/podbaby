CREATE OR REPLACE FUNCTION upsert_channel(my_url text, my_title text, my_description text, my_image text, my_categories text) RETURNS integer AS
$$
DECLARE
attempts INTEGER := 1;
channel_id INTEGER;
BEGIN
LOOP

UPDATE channels SET description=my_description, title=my_title, image=my_image, categories=my_categories WHERE url=my_url RETURNING id INTO channel_id;

IF found THEN
RETURN channel_id;
END IF;

BEGIN
INSERT INTO channels (title, description, image, url, categories) VALUES (my_title, my_description, my_image, my_url, my_categories) RETURNING id INTO channel_id;
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
