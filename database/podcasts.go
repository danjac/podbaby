package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

// PodcastDB manages DB queries to podcasts
type PodcastDB interface {
	SelectAll(int64) ([]models.Podcast, error)
	Create(*models.Podcast) error
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) SelectAll(userID int64) ([]models.Podcast, error) {
	sql := `SELECT DISTINCT p.id, p.title, p.enclosure_url, p.description,
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

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {

	query, args, err := sqlx.Named("SELECT insert_podcast(:channel_id, :title, :description, :enclosure_url, :pub_date)", pc)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)

}
