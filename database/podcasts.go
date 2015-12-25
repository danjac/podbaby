package database

import (
	"github.com/danjac/podbaby/models"
	"github.com/jmoiron/sqlx"
)

type PodcastDB interface {
	SelectAll() ([]models.Podcast, error)
	Create(*models.Podcast) error
}

type defaultPodcastDBImpl struct {
	*sqlx.DB
}

func (db *defaultPodcastDBImpl) SelectAll() ([]models.Podcast, error) {
	sql := `SELECT p.id, p.title, p.enclosure_url, p.description, 
        p.channel_id, c.title AS name, c.image, p.pub_date
        FROM podcasts p 
        JOIN channels c ON c.id = p.channel_id
        ORDER BY pub_date DESC
        LIMIT 30`
	var podcasts []models.Podcast
	err := db.Select(&podcasts, sql)
	return podcasts, err
}

func (db *defaultPodcastDBImpl) Create(pc *models.Podcast) error {
	query, args, err := sqlx.Named(`
        INSERT INTO podcasts (channel_id, title, description, enclosure_url, pub_date) 
        VALUES(:channel_id, :title, :description, :enclosure_url, :pub_date)
        RETURNING id`, pc)

	if err != nil {
		return err
	}

	return db.QueryRow(db.Rebind(query), args...).Scan(&pc.ID)

}
