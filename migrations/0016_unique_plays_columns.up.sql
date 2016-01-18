ALTER TABLE plays ADD CONSTRAINT unique_plays_columns UNIQUE(podcast_id, user_id);
