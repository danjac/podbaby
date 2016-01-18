ALTER TABLE podcasts ADD COLUMN guid TEXT;
ALTER TABLE podcasts ADD CONSTRAINT unique_guids UNIQUE(channel_id, guid);
ALTER TABLE podcasts DROP CONSTRAINT unique_podcast_columns;
