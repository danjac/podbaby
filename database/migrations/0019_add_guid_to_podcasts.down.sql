ALTER TABLE podcasts DROP CONSTRAINT unique_guids;
ALTER TABLE podcasts DROP COLUMN guid;
ALTER TABLE podcasts ADD CONSTRAINT unique_podcast_columns UNIQUE(channel_id, enclosure_url);
