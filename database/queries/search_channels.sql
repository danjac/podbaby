SELECT c.id, c.title, c.description, c.url, c.image, c.website
FROM channels c, plainto_tsquery($1) as q
WHERE (c.tsv @@ q)
ORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20

