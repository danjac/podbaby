package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

const maxSearchRows = 20

// PodcastDB manages DB queries to podcasts
type PodcastDB interface {
	GetByID(int64) (*models.Podcast, error)
	SelectAll(int64) (*models.PodcastList, error)
	SelectSubscribed(int64, int64) (*models.PodcastList, error)
	SelectByChannelID(int64, int64) (*models.PodcastList, error)
	SelectBookmarked(int64, int64) (*models.PodcastList, error)
	SelectPlayed(int64, int64) (*models.PodcastList, error)
	Search(string) ([]models.Podcast, error)
	SearchBookmarked(string, int64) ([]models.Podcast, error)
	SearchByChannelID(string, int64) ([]models.Podcast, error)
	Create(*models.Podcast) error
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) GetByID(id int64) (*models.Podcast, error) {

	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE p.id=$1`

	podcast := &models.Podcast{}
	err := db.Get(podcast, sql, id)
	return podcast, err
}

func (db *defaultPodcastDBImpl) Search(query string) ([]models.Podcast, error) {

	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image
    FROM podcasts p, plainto_tsquery($1) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC LIMIT $2`

	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, sql, query, maxSearchRows)

}

func (db *defaultPodcastDBImpl) SearchByChannelID(query string, channelID int64) ([]models.Podcast, error) {

	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image
    FROM podcasts p, plainto_tsquery($2) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id AND c.id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`

	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, sql, channelID, query, maxSearchRows)

}

func (db *defaultPodcastDBImpl) SearchBookmarked(query string, userID int64) ([]models.Podcast, error) {

	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image
    FROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b
    WHERE (p.tsv @@ q OR c.tsv @@ q) AND p.channel_id = c.id AND b.podcast_id = p.id AND b.user_id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`

	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, sql, userID, query, maxSearchRows)

}

func (db *defaultPodcastDBImpl) SelectPlayed(userID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    WHERE pl.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, p.pub_date, c.title AS name, c.image
    FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    JOIN channels c ON c.id = p.channel_id
    WHERE pl.user_id=$1
    GROUP BY p.id, c.title, c.image, pl.created_at
    ORDER BY pl.created_at DESC
    OFFSET $2 LIMIT $3`

	err := db.Select(
		&result.Podcasts,
		sql,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err

}

func (db *defaultPodcastDBImpl) SelectAll(page int64) (*models.PodcastList, error) {

	sql := "SELECT COUNT(id) FROM podcasts"

	var numRows int64

	if err := db.QueryRow(sql).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    ORDER BY p.pub_date DESC
    OFFSET $1 LIMIT $2`
	err := db.Select(
		&result.Podcasts,
		sql,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectSubscribed(userID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
  JOIN subscriptions s ON s.channel_id = p.channel_id
  WHERE s.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date
    FROM podcasts p
    JOIN subscriptions s ON s.channel_id = p.channel_id
    JOIN channels c ON c.id = p.channel_id
    WHERE s.user_id=$1
    GROUP BY p.id, p.title, c.image, c.title
    ORDER BY p.pub_date DESC
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

	sql := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
  JOIN bookmarks b ON b.podcast_id = p.id
  WHERE b.user_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT p.id, p.title, p.enclosure_url, p.description,
    p.channel_id, c.title AS name, c.image, p.pub_date
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    GROUP BY p.id, p.title, c.title, c.image, b.id
    ORDER BY b.id DESC
    OFFSET $2 LIMIT $3`

	err := db.Select(
		&result.Podcasts,
		sql,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectByChannelID(channelID, page int64) (*models.PodcastList, error) {

	sql := `SELECT COUNT(id) FROM podcasts WHERE channel_id=$1`

	var numRows int64

	if err := db.QueryRow(sql, channelID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	sql = `SELECT id, title, enclosure_url, description, pub_date
    FROM podcasts
    WHERE channel_id=$1
    ORDER BY pub_date DESC
    OFFSET $2 LIMIT $3`

	err := db.Select(
		&result.Podcasts,
		sql,
		channelID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {
	query, args, err := sqlx.Named(`
        SELECT insert_podcast(
            :channel_id, 
            :guid,
            :title, 
            :description, 
            :enclosure_url, 
            :pub_date)`, pc)
	if err != nil {
		return err
	}
	return db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)
}
