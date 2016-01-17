SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    GROUP BY p.id, p.title, c.title, c.image, b.id
    ORDER BY b.id DESC
    OFFSET $2 LIMIT $3
