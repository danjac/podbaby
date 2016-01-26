ALTER TABLE channels DROP COLUMN num_podcasts;

DROP TRIGGER IF EXISTS podcasts_increment_channel_num_podcasts ON podcasts;
DROP TRIGGER IF EXISTS podcasts_decrement_channel_num_podcasts ON podcasts;

DROP FUNCTION increment_num_podcasts_for_channel();
DROP FUNCTION decrement_num_podcasts_for_channel();


