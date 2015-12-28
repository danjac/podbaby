package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// PodcastDB manages DB queries to podcasts
type PodcastDB interface {
	SelectAll(int64, int64) (*models.PodcastList, error)
	SelectByChannelID(int64, int64) ([]models.Podcast, error)
	SelectBookmarked(int64) ([]models.Podcast, error)
	Create(*models.Podcast) error
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) SelectAll(userID int64, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(p.id) FROM podcasts p
  JOIN subscriptions s ON s.channel_id = p.channel_id
  WHERE s.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{}

	result.Page = models.NewPage(page, numRows)

	sql = `SELECT DISTINCT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date,
    EXISTS(SELECT id FROM bookmarks WHERE podcast_id=c.id AND user_id=$1)
      AS is_bookmarked
    FROM podcasts p
    JOIN subscriptions s ON s.channel_id = p.channel_id
    JOIN channels c ON c.id = p.channel_id
    WHERE s.user_id=$1
    ORDER BY pub_date DESC
    OFFSET $2 LIMIT 30`
	err := db.Select(&result.Podcasts, sql, userID, result.Page.Offset)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectBookmarked(userID int64) ([]models.Podcast, error) {
	sql := `SELECT DISTINCT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date,
    EXISTS(SELECT id FROM subscriptions WHERE channel_id=p.channel_id AND user_id=$1)
      AS is_subscribed
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    ORDER BY p.title
    LIMIT 30`
	var podcasts []models.Podcast
	err := db.Select(&podcasts, sql, userID)
	return podcasts, err
}

func (db *defaultPodcastDBImpl) SelectByChannelID(channelID int64, userID int64) ([]models.Podcast, error) {
	sql := `SELECT id, title, enclosure_url, description, pub_date,
    EXISTS(SELECT id FROM bookmarks WHERE podcast_id=podcasts.id AND user_id=$1)
      AS is_bookmarked
    FROM podcasts
    WHERE channel_id=$2
    ORDER BY pub_date DESC
    LIMIT 30`
	var podcasts []models.Podcast
	err := db.Select(&podcasts, sql, userID, channelID)
	return podcasts, err
}

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {
	query, args, err := sqlx.Named("SELECT insert_podcast(:channel_id, :title, :description, :enclosure_url, :pub_date)", pc)
	if err != nil {
		return err
	}
	return db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)
}
