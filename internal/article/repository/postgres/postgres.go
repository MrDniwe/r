package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	"github.com/sirupsen/logrus"
)

func NewRepository(client *sql.DB, l *logrus.Logger) (*ArcticleRepo, error) {
	return &ArcticleRepo{client, l}, nil
}

type ArcticleRepo struct {
	db *sql.DB
	L  *logrus.Logger
}

func (a *ArcticleRepo) GetById(id string) (*models.Article, error) {
	query := `select 
	  uuid, title, lead, body, active_from, views
	from articles
	  where is_visible=true
	  and uuid=$1`
	row := a.db.QueryRow(query, id)
	article := &models.Article{}
	if err := row.Scan(&article.Id, &article.Header, &article.Lead, &article.Text, &article.Date, &article.Views); err != nil {
		return nil, err
	}
	return article, nil
}
