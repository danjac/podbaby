SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image, p.source
FROM podcasts p
JOIN plays pl ON pl.podcast_id = p.id
JOIN channels c ON c.id = p.channel_id
WHERE pl.user_id=$1
GROUP BY p.id, c.title, c.image, pl.created_at
ORDER BY pl.created_at DESC
OFFSET $2 LIMIT $3

