package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

const maxSearchRows = 20

type PodcastReader interface {
	GetByID(DataHandler, int64) (*models.Podcast, error)
	SelectAll(DataHandler, *models.PodcastList, int64) error
	SelectSubscribed(DataHandler, int64, int64) (*models.PodcastList, error)
	SelectByChannel(DataHandler, *models.Channel, int64) (*models.PodcastList, error)
	SelectBookmarked(DataHandler, int64, int64) (*models.PodcastList, error)
	SelectPlayed(DataHandler, int64, int64) (*models.PodcastList, error)
	Search(DataHandler, string) ([]models.Podcast, error)
	SearchBookmarked(DataHandler, string, int64) ([]models.Podcast, error)
	SearchByChannelID(DataHandler, string, int64) ([]models.Podcast, error)
}

type PodcastStore interface {
	PodcastReader
}

type podcastSqlStore struct {
	PodcastReader
}

func newPodcastStore() PodcastStore {

	return &podcastSqlStore{
		PodcastReader: &podcastSqlReader{},
	}

}

type podcastSqlReader struct{}

func (r *podcastSqlReader) GetByID(dh DataHandler, id int64) (*models.Podcast, error) {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, 
        c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE p.id=$1`
	podcast := &models.Podcast{}
	return podcast, sqlx.Get(dh, podcast, q, id)
}

func (r *podcastSqlReader) Search(dh DataHandler, query string) ([]models.Podcast, error) {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($1) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC LIMIT $2`
	var podcasts []models.Podcast
	return podcasts, sqlx.Select(dh, &podcasts, q, query, maxSearchRows)
}

func (r *podcastSqlReader) SearchByChannelID(dh DataHandler, query string, channelID int64) ([]models.Podcast, error) {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
       p.channel_id, p.pub_date, c.title AS name, 
       c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id AND c.id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`
	var podcasts []models.Podcast
	return podcasts, sqlx.Select(dh, &podcasts, q, channelID, query, maxSearchRows)

}

func (r *podcastSqlReader) SearchBookmarked(dh DataHandler, query string, userID int64) ([]models.Podcast, error) {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b
    WHERE (p.tsv @@ q OR c.tsv @@ q) 
        AND p.channel_id = c.id 
        AND b.podcast_id = p.id 
        AND b.user_id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`
	var podcasts []models.Podcast
	return podcasts, sqlx.Select(dh, &podcasts, q, userID, query, maxSearchRows)

}

func (r *podcastSqlReader) SelectPlayed(dh DataHandler, userID, page int64) (*models.PodcastList, error) {

	q := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    WHERE pl.user_id=$1`

	var numRows int64

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPaginator(page, numRows),
	}

	if numRows == 0 {
		return result, nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    JOIN channels c ON c.id = p.channel_id
    WHERE pl.user_id=$1
    GROUP BY p.id, c.title, c.image, pl.created_at
    ORDER BY pl.created_at DESC
    OFFSET $2 LIMIT $3`

	err := sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err

}

func (r *podcastSqlReader) SelectAll(dh DataHandler, result *models.PodcastList, page int64) error {
	var numRows int64

	q := "SELECT reltuples::bigint AS count FROM pg_class WHERE oid = 'public.podcasts'::regclass"

	if err := dh.QueryRowx(q).Scan(&numRows); err != nil {
		return err
	}

	result.Page = models.NewPaginator(page, numRows)

	if numRows == 0 {
		return nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    ORDER BY p.pub_date DESC
    OFFSET $1 LIMIT $2`

	return sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		result.Page.Offset,
		result.Page.PageSize)
}

func (r *podcastSqlReader) SelectSubscribed(dh DataHandler, userID, page int64) (*models.PodcastList, error) {

	q := `
    SELECT SUM(num_podcasts) FROM channels
    WHERE id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)`

	var numRows int64

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		// scan error, as result may be NULL
		if err != nil {
			numRows = 0
		}
	}

	result := &models.PodcastList{
		Page: models.NewPaginator(page, numRows),
	}

	if numRows == 0 {
		return result, nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE c.id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)
    ORDER BY p.pub_date DESC
    OFFSET $2 LIMIT $3`

	err := sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (r *podcastSqlReader) SelectBookmarked(dh DataHandler, userID, page int64) (*models.PodcastList, error) {

	q := `SELECT COUNT(id) FROM bookmarks WHERE user_id=$1`

	var numRows int64

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return nil, err
	}

	result := &models.PodcastList{
		Page: models.NewPaginator(page, numRows),
	}

	if numRows == 0 {
		return result, nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    GROUP BY p.id, p.title, c.title, c.image, b.id
    ORDER BY b.id DESC
    OFFSET $2 LIMIT $3`

	err := sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}

func (r *podcastSqlReader) SelectByChannel(dh DataHandler, channel *models.Channel, page int64) (*models.PodcastList, error) {

	result := &models.PodcastList{
		Page: models.NewPaginator(page, channel.NumPodcasts),
	}

	if channel.NumPodcasts == 0 {
		return result, nil
	}

	q := `
    SELECT id, title, enclosure_url, description, pub_date, source
    FROM podcasts
    WHERE channel_id=$1
    ORDER BY pub_date DESC
    OFFSET $2 LIMIT $3`

	err := sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		channel.ID,
		result.Page.Offset,
		result.Page.PageSize)
	return result, err
}
