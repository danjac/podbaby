ALTER TABLE podcasts ADD COLUMN source TEXT NOT NULL DEFAULT '';

CREATE OR REPLACE FUNCTION insert_podcast (my_channel_id INTEGER, my_guid TEXT, my_title TEXT, my_desc TEXT, my_url TEXT, my_source TEXT, my_date TIMESTAMP) RETURNS INTEGER AS
$$
DECLARE
podcast_id integer;
BEGIN
BEGIN
  INSERT INTO podcasts (channel_id, guid, title, description, enclosure_url, source, pub_date)
  VALUES(my_channel_id, my_guid, my_title, my_desc, my_url, my_source, my_date)
  RETURNING id INTO podcast_id;
EXCEPTION WHEN unique_violation THEN
  SELECT id FROM podcasts WHERE channel_id=my_channel_id AND guid=my_guid INTO podcast_id;
  UPDATE podcasts SET enclosure_url=my_url, title=my_title, source=my_source, description=my_desc, pub_date=my_date WHERE id=podcast_id;
END;
RETURN podcast_id;
END;
$$
LANGUAGE plpgsql;
