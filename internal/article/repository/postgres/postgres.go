package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/internal/server"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewRepository(srv *server.Server) (*ArcticleRepo, error) {
	return &ArcticleRepo{srv}, nil
}

type ArcticleRepo struct {
	Srv *server.Server
}

func (a *ArcticleRepo) GetById(id string) (*models.Article, error) {
	query := `
	select 
		uuid, title, lead, body, active_from, views, image,
		(
			select 
				coalesce(array_to_json(array_agg(com)),'[]') 
			from (
				select 
					uuid, user_uuid, message, created_at,
					(select row_to_json(u) from
						( select 
							uuid, login, email
						from users
						where uuid = c.user_uuid ) u
					) as user
				from comments c
				where
					article_uuid=$1
					and is_visible = true
				order by
					created_at
			) as com
		) as comments
	from articles
	  where is_visible=true
	  and uuid=$1`
	row := a.Srv.Db.QueryRow(query, id)
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
			uuid, title, lead, body, active_from, views, image, '[]' as comments
		from articles
		where
			is_visible = true
		order by active_from desc
		limit $1
		offset $2
	`
	rows, err := a.Srv.Db.Query(query, limit, offset)
	if err != nil {
		nerr := errors.Wrap(err, "Cannot get articles with limit and offset")
		if err, ok := nerr.(e.StackTracer); ok {
			st := err.StackTrace()
			a.Srv.Logger.WithFields(logrus.Fields{
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
	row := a.Srv.Db.QueryRow(query, inPage)
	var total int
	if err := row.Scan(&total); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return 0, e.NotFoundErr
		default:
			nerr := errors.Wrap(err, "Cannot scan row while scanning total article pages")
			if err, ok := nerr.(e.StackTracer); ok {
				a.Srv.Logger.WithFields(logrus.Fields{
					"type":  e.PostgresError,
					"stack": err.StackTrace()[0],
				}).Error(err)
			}
			return 0, e.ServerErr
		}
	}
	return total, nil
}
