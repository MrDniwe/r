package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/pkg/errors"
	"github.com/sirupsen/logrus"
	"html/template"
)

type articleNullable struct {
	Id     string
	Header sql.NullString
	Lead   sql.NullString
	Text   sql.NullString
	Date   pq.NullTime
	Photo  sql.NullString
	Views  sql.NullInt64
}

type scanner interface {
	Scan(...interface{}) error
}

func (a *ArcticleRepo) scanArticle(row scanner) (*models.Article, error) {
	articleNull := articleNullable{}
	if err := row.Scan(&articleNull.Id, &articleNull.Header, &articleNull.Lead, &articleNull.Text, &articleNull.Date, &articleNull.Views, &articleNull.Photo); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.NotFoundErr
		default:
			a.L.WithFields(logrus.Fields{
				"type": "Postgres error",
				"in":   "Scan",
			}).Error(err)
			return nil, errors.ServerErr
		}
		//TODO писать стектрейс, унифицировать типы ошибок
	}
	article := &models.Article{
		Id:      articleNull.Id,
		Visible: true,
		Photo:   articleNull.Photo.String,
		Header:  articleNull.Header.String,
		Lead:    articleNull.Lead.String,
		Text:    template.HTML(articleNull.Text.String),
		Date:    articleNull.Date.Time,
		Views:   int(articleNull.Views.Int64),
	}
	return article, nil
}
