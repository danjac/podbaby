SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image, p.source
FROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b
WHERE (p.tsv @@ q OR c.tsv @@ q) 
    AND p.channel_id = c.id 
    AND b.podcast_id = p.id 
    AND b.user_id=$1
ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3


