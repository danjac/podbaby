ALTER TABLE channels ADD COLUMN num_podcasts INTEGER DEFAULT 0;

CREATE OR REPLACE FUNCTION increment_num_podcasts_for_channel()
RETURNS TRIGGER AS
$$
   BEGIN
    UPDATE channels SET num_podcasts = num_podcasts + 1 WHERE id=new.channel_id;
   RETURN new;
   END;
$$
LANGUAGE 'plpgsql';

CREATE OR REPLACE FUNCTION decrement_num_podcasts_for_channel()
RETURNS TRIGGER AS
$$
   BEGIN
    UPDATE channels SET num_podcasts = num_podcasts - 1 WHERE id=old.channel_id;
   RETURN old;
   END;
$$
LANGUAGE 'plpgsql';

CREATE TRIGGER podcasts_increment_channel_num_podcasts BEFORE INSERT
  ON podcasts FOR EACH ROW EXECUTE PROCEDURE increment_num_podcasts_for_channel();

CREATE TRIGGER podcasts_decrement_channel_num_podcasts BEFORE DELETE
  ON podcasts FOR EACH ROW EXECUTE PROCEDURE decrement_num_podcasts_for_channel();

UPDATE channels SET num_podcasts = (
    SELECT COUNT(id) FROM podcasts WHERE podcasts.channel_id=channels.id
)
