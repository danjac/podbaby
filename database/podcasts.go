package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// PodcastDB manages DB queries to podcasts
type PodcastDB interface {
	SelectAll(int64) ([]models.Podcast, error)
	GetOrCreate(*models.Podcast) (bool, error)
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) SelectAll(userID int64) ([]models.Podcast, error) {
	sql := `SELECT p.id, p.title, p.enclosure_url, p.description,
        p.channel_id, c.title AS name, c.image, p.pub_date
        FROM podcasts p
        JOIN subscriptions s ON s.channel_id = p.id
        JOIN channels c ON c.id = p.channel_id
        WHERE s.user_id=$1
        ORDER BY pub_date DESC
        LIMIT 30`
	var podcasts []models.Podcast
	err := db.Select(&podcasts, sql, userID)
	return podcasts, err
}

func (db *defaultPodcastDBImpl) GetOrCreate(pc *models.Podcast) (bool, error) {

    query, args, err := sqlx.Named(`
        SELECT * FROM podcasts WHERE enclosure_url=:enclosure_url AND channel_id=:channel_id
    `, pc)
    if err != nil {
        return false, err
    }
    if err := db.Get(pc, query, ...args); err == nil {
        return false, nil
    }
	query, args, err = sqlx.Named(`
        INSERT INTO podcasts (channel_id, title, description, enclosure_url, pub_date)
        VALUES(:channel_id, :title, :description, :enclosure_url, :pub_date)
        RETURNING id`, pc)

	if err != nil {
		return false, err
	}

	return true, db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)

}
