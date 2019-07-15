package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/pkg/errors"
	ers "github.com/pkg/errors"
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
	art, err := a.scanArticle(row)
	if err != nil {
		return nil, err
	}
	return art, nil
}

func (a *ArcticleRepo) GetLastNArticles(n int) ([]*models.Article, error) {
	query := `
	  select
	    uuid, title,, lead, body, active_from, views, image
	  from articles
	    where
	      is_visible=true
	    order by active_from desc
	    limit $1
	`
	rows, err := a.db.Query(query, n)
	if err != nil {
		nerr := ers.Wrap(err, "Cannot get last N articles")
		if err, ok := nerr.(errors.StackTracer); ok {
			st := err.StackTrace()
			a.L.WithFields(logrus.Fields{
				"stack": fmt.Sprintf("%+v", st[0]),
				"type":  errors.ServerError,
			}).Error(err)
		}
		//TODO обернуть в какую-то функцию
		return nil, errors.ServerErr
	}
	defer rows.Close()
	articles := make([]*models.Article, n)
	for i := 0; rows.Next(); i++ {
		articles[i], err = a.scanArticle(rows)
		if err != nil {
			return nil, err
		}
	}
	return articles, nil
}
