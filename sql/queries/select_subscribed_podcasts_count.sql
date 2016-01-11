SELECT COUNT(DISTINCT(id)) FROM podcasts
WHERE channel_id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)


