package store

import (
	"github.com/danjac/podbaby/models"
	"github.com/danjac/podbaby/store/Godeps/_workspace/src/github.com/jmoiron/sqlx"
)

const maxSearchRows = 20

// PodcastReader handles reads from podcast data store
type PodcastReader interface {
	GetByID(DataHandler, *models.Podcast, int) error
	SelectAll(DataHandler, *models.PodcastList, int) error
	SelectSubscribed(DataHandler, *models.PodcastList, int, int) error
	SelectByChannel(DataHandler, *models.PodcastList, *models.Channel, int) error
	SelectBookmarked(DataHandler, *models.PodcastList, int, int) error
	SelectPlayed(DataHandler, *models.PodcastList, int, int) error
	Search(DataHandler, *models.PodcastList, string, int) error
	SearchBookmarked(DataHandler, *[]models.Podcast, string, int) error
	SearchByChannelID(DataHandler, *[]models.Podcast, string, int) error
}

// PodcastStore handles interactions with podcast data store
type PodcastStore interface {
	PodcastReader
}

type podcastSQLStore struct {
	PodcastReader
}

func newPodcastStore() PodcastStore {

	return &podcastSQLStore{
		PodcastReader: &podcastSQLReader{},
	}

}

type podcastSQLReader struct{}

func (r *podcastSQLReader) GetByID(dh DataHandler, podcast *models.Podcast, id int) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, 
        c.image, p.pub_date, p.source
    FROM podcasts p
    JOIN channels c ON c.id = p.channel_id
    WHERE p.id=$1`
	return handleError(sqlx.Get(dh, podcast, q, id), q)
}

func (r *podcastSQLReader) Search(dh DataHandler, result *models.PodcastList, query string, page int) error {
	q := `SELECT COUNT(p.id) FROM podcasts p, plainto_tsquery($1) as q WHERE (p.tsv @@ q)`

	var numRows int

	if err := dh.QueryRowx(q, query).Scan(&numRows); err != nil {
		return handleError(err, q)
	}

	result.Page = models.NewPaginator(page, numRows)

	if numRows == 0 {
		return nil
	}

	q = `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($1) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($1)) DESC, p.pub_date DESC
    OFFSET $2 LIMIT $3`

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		query,
		result.Page.Offset,
		result.Page.PageSize), q)

}

func (r *podcastSQLReader) SearchByChannelID(dh DataHandler, podcasts *[]models.Podcast, query string, channelID int) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
       p.channel_id, p.pub_date, c.title AS name, 
       c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c
    WHERE (p.tsv @@ q) AND p.channel_id = c.id AND c.id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC, p.pub_date DESC LIMIT $3`
	return handleError(sqlx.Select(dh, podcasts, q, channelID, query, maxSearchRows), q)

}

func (r *podcastSQLReader) SearchBookmarked(dh DataHandler, podcasts *[]models.Podcast, query string, userID int) error {
	q := `
    SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, p.pub_date, c.title AS name, c.image, p.source
    FROM podcasts p, plainto_tsquery($2) as q, channels c, bookmarks b
    WHERE (p.tsv @@ q OR c.tsv @@ q) 
        AND p.channel_id = c.id 
        AND b.podcast_id = p.id 
        AND b.user_id=$1
    ORDER BY ts_rank_cd(p.tsv, plainto_tsquery($2)) DESC, p.pub_date DESC LIMIT $3`
	return sqlx.Select(dh, podcasts, q, userID, query, maxSearchRows)

}

func (r *podcastSQLReader) SelectPlayed(dh DataHandler, result *models.PodcastList, userID, page int) error {

	q := `SELECT COUNT(DISTINCT(p.id)) FROM podcasts p
    JOIN plays pl ON pl.podcast_id = p.id
    WHERE pl.user_id=$1`

	var numRows int

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return handleError(err, q)
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

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize), q)

}

func (r *podcastSQLReader) SelectAll(dh DataHandler, result *models.PodcastList, page int) error {
	var numRows int

	q := "SELECT reltuples::bigint AS count FROM pg_class WHERE oid = 'public.podcasts'::regclass"

	if err := dh.QueryRowx(q).Scan(&numRows); err != nil {
		return handleError(err, q)
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

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		result.Page.Offset,
		result.Page.PageSize), q)
}

func (r *podcastSQLReader) SelectSubscribed(dh DataHandler, result *models.PodcastList, userID, page int) error {

	q := `
    SELECT SUM(num_podcasts) FROM channels
    WHERE id IN (SELECT channel_id FROM subscriptions WHERE user_id=$1)`

	var numRows int

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

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize), q)
}

func (r *podcastSQLReader) SelectBookmarked(dh DataHandler, result *models.PodcastList, userID, page int) error {

	q := `SELECT COUNT(id) FROM bookmarks WHERE user_id=$1`

	var numRows int

	if err := dh.QueryRowx(q, userID).Scan(&numRows); err != nil {
		return handleError(err, q)
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

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		userID,
		result.Page.Offset,
		result.Page.PageSize), q)
}

func (r *podcastSQLReader) SelectByChannel(dh DataHandler, result *models.PodcastList, channel *models.Channel, page int) error {

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

	return handleError(sqlx.Select(
		dh,
		&result.Podcasts,
		q,
		channel.ID,
		result.Page.Offset,
		result.Page.PageSize), q)
}
