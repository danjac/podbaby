ALTER TABLE podcasts ADD CONSTRAINT unique_podcast_columns UNIQUE(channel_id, enclosure_url);
