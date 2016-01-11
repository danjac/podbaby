SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c
WHERE url=$1

