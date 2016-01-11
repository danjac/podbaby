package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/sql"
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

	q, _ := sql.Queries.Get("get_podcast_by_id.sql")
	podcast := &models.Podcast{}
	err := db.Get(podcast, q, id)
	return podcast, err
}

func (db *defaultPodcastDBImpl) Search(query string) ([]models.Podcast, error) {

	q, _ := sql.Queries.Get("search_podcasts.sql")
	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, q, query, maxSearchRows)
}

func (db *defaultPodcastDBImpl) SearchByChannelID(query string, channelID int64) ([]models.Podcast, error) {

	q, _ := sql.Queries.Get("search_podcasts_by_channel_id.sql")
	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, q, channelID, query, maxSearchRows)

}

func (db *defaultPodcastDBImpl) SearchBookmarked(query string, userID int64) ([]models.Podcast, error) {

	q, _ := sql.Queries.Get("search_bookmarked_podcasts.sql")
	var podcasts []models.Podcast
	return podcasts, db.Select(&podcasts, q, userID, query, maxSearchRows)

}

func (db *defaultPodcastDBImpl) SelectPlayed(userID, page int64) (*models.PodcastList, error) {

	q, _ := sql.Queries.Get("select_played_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRow(q, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = sql.Queries.Get("select_played_podcasts.sql")

	err := db.Select(
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err

}

func (db *defaultPodcastDBImpl) SelectAll(page int64) (*models.PodcastList, error) {

	q, _ := sql.Queries.Get("select_all_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRow(q).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = sql.Queries.Get("select_all_podcasts.sql")

	err := db.Select(
		&result.Podcasts,
		q,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectSubscribed(userID, page int64) (*models.PodcastList, error) {

	q, _ := sql.Queries.Get("select_subscribed_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRow(q, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = sql.Queries.Get("select_subscribed_podcasts.sql")

	err := db.Select(
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectBookmarked(userID, page int64) (*models.PodcastList, error) {

	q, _ := sql.Queries.Get("select_bookmarked_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRow(q, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = sql.Queries.Get("select_bookmarked_podcasts.sql")

	err := db.Select(
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) SelectByChannelID(channelID, page int64) (*models.PodcastList, error) {

	q, _ := sql.Queries.Get("select_podcasts_by_channel_id_count.sql")

	var numRows int64

	if err := db.QueryRow(q, channelID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = sql.Queries.Get("select_podcasts_by_channel_id.sql")

	err := db.Select(
		&result.Podcasts,
		q,
		channelID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {
	q, _ := sql.Queries.Get("insert_podcast.sql")
	q, args, err := sqlx.Named(q, pc)
	if err != nil {
		return err
	}
	return db.QueryRow(db.Rebind(q), args...).Scan(&pc.ID)
}
