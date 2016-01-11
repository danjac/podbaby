SELECT c.id, c.title, c.description, c.image, c.url, c.website
FROM channels c
JOIN subscriptions s ON s.channel_id = c.id
WHERE s.user_id=$1 AND title IS NOT NULL AND title != ''
GROUP BY c.id
ORDER BY title

