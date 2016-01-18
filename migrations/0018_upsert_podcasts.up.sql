CREATE OR REPLACE FUNCTION insert_podcast (my_channel_id INTEGER, my_title TEXT, my_desc TEXT, my_url TEXT, my_date TIMESTAMP) RETURNS INTEGER AS
$$
DECLARE
podcast_id integer;
BEGIN
BEGIN
  INSERT INTO podcasts (channel_id, title, description, enclosure_url, pub_date)
  VALUES(my_channel_id, my_title, my_desc, my_url, my_date)
  RETURNING id INTO podcast_id;
EXCEPTION WHEN unique_violation THEN
  SELECT id FROM podcasts WHERE channel_id=my_channel_id AND enclosure_url=my_url INTO podcast_id;
  UPDATE podcasts SET title=my_title, description=my_desc, pub_date=my_date WHERE id=podcast_id;
END;
RETURN podcast_id;
END;
$$
LANGUAGE plpgsql;
