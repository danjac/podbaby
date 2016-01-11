SELECT id, title, enclosure_url, description, pub_date, source
FROM podcasts
WHERE channel_id=$1
ORDER BY pub_date DESC
OFFSET $2 LIMIT $3


