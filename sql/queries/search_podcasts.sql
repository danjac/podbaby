SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image, p.source
FROM podcasts p, plainto_tsquery($1) as q, channels c
WHERE (p.tsv @@ q) AND p.channel_id = c.id
ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC LIMIT $2


