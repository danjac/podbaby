package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// PodcastDB manages DB queries to podcasts
type PodcastDB interface {
	SelectSubscribed(int64, int64) (*models.PodcastList, error)
	SelectByChannelID(int64, int64, int64) (*models.PodcastList, error)
	SelectBookmarked(int64, int64) (*models.PodcastList, error)
	Search(string, int64) ([]models.Podcast, error)
	Create(*models.Podcast) error
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) Search(query string, userID int64) ([]models.Podcast, error) {

	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image,
    EXISTS(SELECT id FROM bookmarks WHERE podcast_id=p.id AND user_id=$1)
      AS is_bookmarked,
    EXISTS(SELECT id FROM subscriptions WHERE channel_id=p.channel_id AND user_id=$1)
      AS is_subscribed
    FROM podcasts p, plainto_tsquery($2) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT 20`

	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, sql, userID, query)

}

func (db *defaultPodcastDBImpl) SelectSubscribed(userID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(p.id) FROM podcasts p
  JOIN subscriptions s ON s.channel_id = p.channel_id
  WHERE s.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT DISTINCT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date,
    EXISTS(SELECT id FROM bookmarks WHERE podcast_id=p.id AND user_id=$1)
      AS is_bookmarked,
    1 AS is_subscribed
    FROM podcasts p
    JOIN subscriptions s ON s.channel_id = p.channel_id
    JOIN channels c ON c.id = p.channel_id
    WHERE s.user_id=$1
    ORDER BY pub_date DESC
    OFFSET $2 LIMIT $3`
	err := db.Select(
		&result.Podcasts,
		sql,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectBookmarked(userID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(p.id) FROM podcasts p
  JOIN bookmarks b ON b.podcast_id = p.id
  WHERE b.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT DISTINCT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date,
    EXISTS(SELECT id FROM subscriptions WHERE channel_id=p.channel_id AND user_id=$1)
      AS is_subscribed,
    1 AS is_bookmarked
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    ORDER BY p.title
    OFFSET $2 LIMIT $3`

	err := db.Select(
		&result.Podcasts,
		sql,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectByChannelID(channelID, userID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(id) FROM podcasts WHERE channel_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, channelID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT id, title, enclosure_url, description, pub_date,
    EXISTS(SELECT id FROM bookmarks WHERE podcast_id=podcasts.id AND user_id=$1)
      AS is_bookmarked
    FROM podcasts
    WHERE channel_id=$2
    ORDER BY pub_date DESC
    OFFSET $3 LIMIT $4`

	err := db.Select(
		&result.Podcasts,
		sql,
		userID,
		channelID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {
	query, args, err := sqlx.Named("SELECT insert_podcast(:channel_id, :title, :description, :enclosure_url, :pub_date)", pc)
	if err != nil {
		return err
	}
	return db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)
}
