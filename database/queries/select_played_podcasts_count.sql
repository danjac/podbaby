SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
JOIN plays pl ON pl.podcast_id = p.id
WHERE pl.user_id=$1


