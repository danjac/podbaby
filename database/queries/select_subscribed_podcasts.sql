SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE c.id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)
    ORDER BY p.pub_date DESC
    OFFSET $2 LIMIT $3
