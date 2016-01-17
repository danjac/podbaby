package database

// Queries is a *GenPurse.
	var Queries = &GenPurse{
	files: map[string]string{
		
			"add_play.sql": "SELECT add_play($1, $2)\n",
		
			"delete_bookmark.sql": "DELETE FROM bookmarks WHERE podcast_id=$1 AND user_id=$2\n",
		
			"delete_plays_by_user_id.sql": "DELETE FROM plays WHERE user_id=$1\n",
		
			"delete_subscription.sql": "DELETE FROM subscriptions WHERE channel_id=$1 AND user_id=$2\n",
		
			"delete_user.sql": "DELETE FROM users WHERE id=$1\n",
		
			"get_channel_by_id.sql": "SELECT c.id, c.title, c.description, c.url, c.image, c.website\nFROM channels c\nWHERE id=$1\n\n",
		
			"get_channel_by_url.sql": "SELECT c.id, c.title, c.description, c.url, c.image, c.website\nFROM channels c\nWHERE url=$1\n\n",
		
			"get_podcast_by_id.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, c.title AS name, c.image, p.pub_date, p.source\nFROM podcasts p\nJOIN channels c ON c.id = p.channel_id\nWHERE p.id=$1\n\n\n",
		
			"get_user_by_id.sql": "SELECT * FROM users WHERE id=$1\n",
		
			"get_user_by_name_or_email.sql": "SELECT * FROM users WHERE email=$1 or name=$1\n",
		
			"insert_bookmark.sql": "INSERT INTO bookmarks(podcast_id, user_id) VALUES($1, $2)\n",
		
			"insert_podcast.sql": "SELECT insert_podcast(\n    :channel_id, \n    :guid,\n    :title, \n    :description, \n    :enclosure_url, \n    :source,\n    :pub_date\n)\n\n",
		
			"insert_subscription.sql": "INSERT INTO subscriptions(channel_id, user_id) VALUES($1, $2)\n",
		
			"insert_user.sql": "INSERT INTO users(name, email, password)\n    VALUES (:name, :email, :password) RETURNING id\n",
		
			"search_bookmarked_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, p.pub_date, c.title AS name, c.image, p.source\nFROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b\nWHERE (p.tsv @@ q OR c.tsv @@ q) \n    AND p.channel_id = c.id \n    AND b.podcast_id = p.id \n    AND b.user_id=$1\nORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3\n\n\n",
		
			"search_channels.sql": "SELECT c.id, c.title, c.description, c.url, c.image, c.website\nFROM channels c, plainto_tsquery($1) as q\nWHERE (c.tsv @@ q)\nORDER BY ts_rank_cd(c.tsv, plainto_tsquery($1)) DESC LIMIT 20\n\n",
		
			"search_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, p.pub_date, c.title AS name, c.image, p.source\nFROM podcasts p, plainto_tsquery($1) as q, channels c\nWHERE (p.tsv @@ q) AND p.channel_id = c.id\nORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC LIMIT $2\n\n\n",
		
			"search_podcasts_by_channel_id.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, p.pub_date, c.title AS name, c.image, p.source\nFROM podcasts p, plainto_tsquery($2) as q, channels c\nWHERE (p.tsv @@ q) AND p.channel_id = c.id AND c.id=$1\nORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3\n\n\n",
		
			"select_all_channels.sql": "SELECT id, title, description, categories, url, image, website FROM channels\n",
		
			"select_all_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, c.title AS name, c.image, p.pub_date, p.source\nFROM podcasts p\nJOIN channels c ON c.id = p.channel_id\nORDER BY p.pub_date DESC\nOFFSET $1 LIMIT $2\n\n",
		
			"select_all_podcasts_count.sql": "SELECT COUNT(id) FROM podcasts\n",
		
			"select_bookmarked_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, c.title AS name, c.image, p.pub_date, p.source\n    FROM podcasts p\n    JOIN channels c ON c.id = p.channel_id\n    JOIN bookmarks b ON b.podcast_id = p.id\n    WHERE b.user_id=$1\n    GROUP BY p.id, p.title, c.title, c.image, b.id\n    ORDER BY b.id DESC\n    OFFSET $2 LIMIT $3\n",
		
			"select_bookmarked_podcasts_count.sql": "SELECT COUNT(DISTINCT(p.id)) FROM podcasts p\n  JOIN bookmarks b ON b.podcast_id = p.id\n  WHERE b.user_id=$1\n\n\n",
		
			"select_bookmarks.sql": "SELECT podcast_id FROM bookmarks WHERE user_id=$1\n",
		
			"select_played_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, p.pub_date, c.title AS name, c.image, p.source\nFROM podcasts p\nJOIN plays pl ON pl.podcast_id = p.id\nJOIN channels c ON c.id = p.channel_id\nWHERE pl.user_id=$1\nGROUP BY p.id, c.title, c.image, pl.created_at\nORDER BY pl.created_at DESC\nOFFSET $2 LIMIT $3\n\n",
		
			"select_played_podcasts_count.sql": "SELECT COUNT(DISTINCT(p.id)) FROM podcasts p\nJOIN plays pl ON pl.podcast_id = p.id\nWHERE pl.user_id=$1\n\n\n",
		
			"select_podcasts_by_channel_id.sql": "SELECT id, title, enclosure_url, description, pub_date, source\nFROM podcasts\nWHERE channel_id=$1\nORDER BY pub_date DESC\nOFFSET $2 LIMIT $3\n\n\n",
		
			"select_podcasts_by_channel_id_count.sql": "SELECT COUNT(id) FROM podcasts WHERE channel_id=$1\n",
		
			"select_subscribed_channels.sql": "SELECT c.id, c.title, c.description, c.image, c.url, c.website\nFROM channels c\nJOIN subscriptions s ON s.channel_id = c.id\nWHERE s.user_id=$1 AND title IS NOT NULL AND title != ''\nGROUP BY c.id\nORDER BY title\n\n",
		
			"select_subscribed_podcasts.sql": "SELECT p.id, p.title, p.enclosure_url, p.description,\n    p.channel_id, c.title AS name, c.image, p.pub_date, p.source\n    FROM podcasts p\n    JOIN channels c ON c.id = p.channel_id\n    WHERE c.id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)\n    ORDER BY p.pub_date DESC\n    OFFSET $2 LIMIT $3\n",
		
			"select_subscribed_podcasts_count.sql": "SELECT COUNT(DISTINCT(id)) FROM podcasts\nWHERE channel_id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)\n\n\n",
		
			"select_subscriptions_by_user_id.sql": "SELECT channel_id FROM subscriptions WHERE user_id=$1\n",
		
			"update_user_email.sql": "UPDATE users SET email=$1 WHERE id=$2\n",
		
			"update_user_password.sql": "UPDATE users SET password=$1 WHERE id=$2\n",
		
			"upsert_channel.sql": "SELECT upsert_channel (\n    :url, \n    :title, \n    :description, \n    :image, \n    :categories, \n    :website\n)\n\n\n",
		
			"user_email_exists.sql": "SELECT COUNT(id) FROM users WHERE email ILIKE $1\n",
		
			"user_email_exists_with_id.sql": "SELECT COUNT(id) FROM users WHERE email ILIKE $1 AND id != $2\n",
		
			"user_name_exists.sql": "SELECT COUNT(id) FROM users WHERE name=$1\n",
		
	},
}
