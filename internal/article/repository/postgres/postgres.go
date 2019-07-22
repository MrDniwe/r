package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/pkg/errors"
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

func (a *ArcticleRepo) GetLastList(limit, offset int) ([]*models.Article, error) {
	if limit > 100 {
		limit = 100
	}
	query := `
		select
			uuid, title, lead, body, active_from, views, image
		from articles
		where
			is_visible = true
		order by active_from desc
		limit $1
		offset $2
	`
	rows, err := a.db.Query(query, limit, offset)
	if err != nil {
		nerr := errors.Wrap(err, "Cannot get articles with limit and offset")
		if err, ok := nerr.(e.StackTracer); ok {
			st := err.StackTrace()
			a.L.WithFields(logrus.Fields{
				"stack": fmt.Sprintf("%+v", st[0]),
				"type":  e.PostgresError,
			}).Error(err)
		}
		return nil, e.ServerErr
	}
	defer rows.Close()
	articles := make([]*models.Article, 0, limit)
	for rows.Next() {
		art, err := a.scanArticle(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}

func (a *ArcticleRepo) PagesCount(inPage int) (int, error) {
	if inPage < 1 {
		return 0, e.BadRequestErr
	}
	query := `select ceil( count(*)/$1 ) as total from articles where is_visible = true`
	row := a.db.QueryRow(query, inPage)
	var total int
	if err := row.Scan(&total); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return 0, e.NotFoundErr
		default:
			nerr := errors.Wrap(err, "Cannot scan row while scanning total article pages")
			if err, ok := nerr.(e.StackTracer); ok {
				a.L.WithFields(logrus.Fields{
					"type":  e.PostgresError,
					"stack": err.StackTrace()[0],
				}).Error(err)
			}
			return 0, e.ServerErr
		}
	}
	return total, nil
}
