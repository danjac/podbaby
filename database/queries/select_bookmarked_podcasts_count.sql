SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
  JOIN bookmarks b ON b.podcast_id = p.id
  WHERE b.user_id=$1


