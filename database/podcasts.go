package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
	"github.com/smotes/purse"
)

const maxSearchRows = 20

type PodcastWriter interface {
	Create(*models.Podcast) error
}

type PodcastReader interface {
	GetByID(int64) (*models.Podcast, error)
	SelectAll(int64) (*models.PodcastList, error)
	SelectSubscribed(int64, int64) (*models.PodcastList, error)
	SelectByChannelID(int64, int64) (*models.PodcastList, error)
	SelectBookmarked(int64, int64) (*models.PodcastList, error)
	SelectPlayed(int64, int64) (*models.PodcastList, error)
	Search(string) ([]models.Podcast, error)
	SearchBookmarked(string, int64) ([]models.Podcast, error)
	SearchByChannelID(string, int64) ([]models.Podcast, error)
}

type PodcastDB struct {
	PodcastReader
	PodcastWriter
}

func newPodcastDB(db sqlx.Ext, ps purse.Purse) *PodcastDB {

	return &PodcastDB{
		PodcastWriter: &PodcastDBWriter{db, ps},
		PodcastReader: &PodcastDBReader{db, ps},
	}

}

type PodcastDBWriter struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *PodcastDBWriter) Create(pc *models.Podcast) error {
	q, _ := db.ps.Get("insert_podcast.sql")
	q, args, err := sqlx.Named(q, pc)
	if err != nil {
		return sqlErr(err, q)
	}
	return sqlErr(db.QueryRowx(db.Rebind(q), args...).Scan(&pc.ID), q)
}

type PodcastDBReader struct {
	sqlx.Ext
	ps purse.Purse
}

func (db *PodcastDBReader) GetByID(id int64) (*models.Podcast, error) {
	q, _ := db.ps.Get("get_podcast_by_id.sql")
	podcast := &models.Podcast{}
	err := sqlx.Get(db, podcast, q, id)
	return podcast, sqlErr(err, q)
}

func (db *PodcastDBReader) Search(query string) ([]models.Podcast, error) {

	q, _ := db.ps.Get("search_podcasts.sql")
	var podcasts []models.Podcast
	return podcasts, sqlErr(sqlx.Select(db, &podcasts, q, query, maxSearchRows), q)
}

func (db *PodcastDBReader) SearchByChannelID(query string, channelID int64) ([]models.Podcast, error) {

	q, _ := db.ps.Get("search_podcasts_by_channel_id.sql")
	var podcasts []models.Podcast
	return podcasts, sqlErr(sqlx.Select(db, &podcasts, q, channelID, query, maxSearchRows), q)

}

func (db *PodcastDBReader) SearchBookmarked(query string, userID int64) ([]models.Podcast, error) {

	q, _ := db.ps.Get("search_bookmarked_podcasts.sql")
	var podcasts []models.Podcast
	return podcasts, sqlErr(sqlx.Select(db, &podcasts, q, userID, query, maxSearchRows), q)

}

func (db *PodcastDBReader) SelectPlayed(userID, page int64) (*models.PodcastList, error) {

	q, _ := db.ps.Get("select_played_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return nil, sqlErr(err, q)
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = db.ps.Get("select_played_podcasts.sql")

	err := sqlx.Select(
		db,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, sqlErr(err, q)

}

func (db *PodcastDBReader) SelectAll(page int64) (*models.PodcastList, error) {

	q, _ := db.ps.Get("select_all_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRowx(q).Scan(&numRows); err != nil {
		return nil, sqlErr(err, q)
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = db.ps.Get("select_all_podcasts.sql")

	err := sqlx.Select(
		db,
		&result.Podcasts,
		q,
		result.Page.Offset,
		result.Page.PageSize)
	return result, sqlErr(err, q)
}

func (db *PodcastDBReader) SelectSubscribed(userID, page int64) (*models.PodcastList, error) {

	q, _ := db.ps.Get("select_subscribed_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return nil, sqlErr(err, q)
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = db.ps.Get("select_subscribed_podcasts.sql")

	err := sqlx.Select(
		db,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, sqlErr(err, q)
}

func (db *PodcastDBReader) SelectBookmarked(userID, page int64) (*models.PodcastList, error) {

	q, _ := db.ps.Get("select_bookmarked_podcasts_count.sql")

	var numRows int64

	if err := db.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return nil, sqlErr(err, q)
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = db.ps.Get("select_bookmarked_podcasts.sql")

	err := sqlx.Select(
		db,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, sqlErr(err, q)
}

func (db *PodcastDBReader) SelectByChannelID(channelID, page int64) (*models.PodcastList, error) {

	q, _ := db.ps.Get("select_podcasts_by_channel_id_count.sql")

	var numRows int64

	if err := db.QueryRowx(q, channelID).Scan(&numRows); err != nil {
		return nil, sqlErr(err, q)
	}

	result := &models.PodcastList{
		Page: models.NewPage(page, numRows),
	}

	q, _ = db.ps.Get("select_podcasts_by_channel_id.sql")

	err := sqlx.Select(
		db,
		&result.Podcasts,
		q,
		channelID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, sqlErr(err, q)
}
