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
	  uuid, title, lead, body, active_from, views, image
	from articles
	  where is_visible=true
	  and uuid=$1`
	row := a.db.QueryRow(query, id)
	return scanArticle(row)
}

func (a *ArcticleRepo) GetLastNArticles(n int) ([]*models.Article, error) {
	query := `
	  select
	    uuid, title, lead, body, active_from, views, image
	  from articles
	    where
	      is_visible=true
	    order by active_from desc
	    limit $1
	`
	rows, err := a.db.Query(query, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := make([]*models.Article, n)
	for i := 0; rows.Next(); i++ {
		articles[i], err = scanArticle(rows)
		if err != nil {
			return nil, err
		}
	}
	return articles, nil
}
