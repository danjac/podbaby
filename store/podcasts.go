package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

const maxSearchRows = 20

type PodcastReader interface {
	GetByID(DataHandler, *models.Podcast, int64) error
	SelectAll(DataHandler, *models.PodcastList, int64) error
	SelectSubscribed(DataHandler, *models.PodcastList, int64, int64) error
	SelectByChannel(DataHandler, *models.PodcastList, *models.Channel, int64) error
	SelectBookmarked(DataHandler, *models.PodcastList, int64, int64) error
	SelectPlayed(DataHandler, *models.PodcastList, int64, int64) error
	Search(DataHandler, *[]models.Podcast, string) error
	SearchBookmarked(DataHandler, *[]models.Podcast, string, int64) error
	SearchByChannelID(DataHandler, *[]models.Podcast, string, int64) error
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

func (r *podcastSqlReader) GetByID(dh DataHandler, podcast *models.Podcast, id int64) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, 
        c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE p.id=$1`
	return sqlx.Get(dh, podcast, q, id)
}

func (r *podcastSqlReader) Search(dh DataHandler, podcasts *[]models.Podcast, query string) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($1) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC LIMIT $2`
	return sqlx.Select(dh, podcasts, q, query, maxSearchRows)
}

func (r *podcastSqlReader) SearchByChannelID(dh DataHandler, podcasts *[]models.Podcast, query string, channelID int64) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
       p.channel_id, p.pub_date, c.title AS name, 
       c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id AND c.id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`
	return sqlx.Select(dh, podcasts, q, channelID, query, maxSearchRows)

}

func (r *podcastSqlReader) SearchBookmarked(dh DataHandler, podcasts *[]models.Podcast, query string, userID int64) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b
    WHERE (p.tsv @@ q OR c.tsv @@ q) 
        AND p.channel_id = c.id 
        AND b.podcast_id = p.id 
        AND b.user_id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC LIMIT $3`
	return sqlx.Select(dh, podcasts, q, userID, query, maxSearchRows)

}

func (r *podcastSqlReader) SelectPlayed(dh DataHandler, result *models.PodcastList, userID, page int64) error {

	q := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    WHERE pl.user_id=$1`

	var numRows int64

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return err
	}

	result.Page = models.NewPaginator(page, numRows)

	if numRows == 0 {
		return nil
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

	return sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)

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

func (r *podcastSqlReader) SelectSubscribed(dh DataHandler, result *models.PodcastList, userID, page int64) error {

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

	result.Page = models.NewPaginator(page, numRows)

	if numRows == 0 {
		return nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE c.id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)
    ORDER BY p.pub_date DESC
    OFFSET $2 LIMIT $3`

	return sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
}

func (r *podcastSqlReader) SelectBookmarked(dh DataHandler, result *models.PodcastList, userID, page int64) error {

	q := `SELECT COUNT(id) FROM bookmarks WHERE user_id=$1`

	var numRows int64

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
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
    JOIN bookmarks b ON b.podcast_id = p.id
    WHERE b.user_id=$1
    GROUP BY p.id, p.title, c.title, c.image, b.id
    ORDER BY b.id DESC
    OFFSET $2 LIMIT $3`

	return sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize)
}

func (r *podcastSqlReader) SelectByChannel(dh DataHandler, result *models.PodcastList, channel *models.Channel, page int64) error {

	result.Page = models.NewPaginator(page, channel.NumPodcasts)

	if channel.NumPodcasts == 0 {
		return nil
	}

	q := `
    SELECT id, title, enclosure_url, description, pub_date, source
    FROM podcasts
    WHERE channel_id=$1
    ORDER BY pub_date DESC
    OFFSET $2 LIMIT $3`

	return sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		channel.ID,
		result.Page.Offset,
		result.Page.PageSize)
}
