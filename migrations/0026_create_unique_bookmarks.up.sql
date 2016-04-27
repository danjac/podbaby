DELETE FROM bookmarks
WHERE id IN (SELECT id
              FROM (SELECT id,
                             ROW_NUMBER() OVER (partition BY podcast_id, user_id ORDER BY id DESC) AS rnum
                     FROM bookmarks) t
              WHERE t.rnum > 1);

CREATE UNIQUE INDEX bookmarks_unique_idx ON bookmarks (podcast_id, user_id);
