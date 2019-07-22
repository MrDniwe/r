package repository

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"html/template"
)

type articleNullable struct {
	Id          string
	Header      sql.NullString
	Lead        sql.NullString
	Text        sql.NullString
	Date        pq.NullTime
	Photo       sql.NullString
	Views       sql.NullInt64
	CommentsAgg sql.NullString
}

type scanner interface {
	Scan(...interface{}) error
}

func (a *ArcticleRepo) scanArticle(row scanner) (*models.Article, error) {
	articleNull := articleNullable{}
	if err := row.Scan(&articleNull.Id, &articleNull.Header, &articleNull.Lead, &articleNull.Text, &articleNull.Date, &articleNull.Views, &articleNull.Photo, &articleNull.CommentsAgg); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, e.NotFoundErr
		default:
			nerr := errors.Wrap(err, "Cannot scan row while scanning articles")
			if err, ok := nerr.(e.StackTracer); ok {
				a.L.WithFields(logrus.Fields{
					"type":  e.PostgresError,
					"stack": err.StackTrace()[0],
				}).Error(err)
			}
			return nil, e.ServerErr
		}
	}
	comments := models.Comments{}
	if err := json.Unmarshal([]byte(articleNull.CommentsAgg.String), &comments); err != nil {
		a.L.WithFields(logrus.Fields{
			"type":    e.UnmarshalError,
			"details": err,
		})
		return nil, e.ServerErr
	}
	article := &models.Article{
		Id:       articleNull.Id,
		Visible:  true,
		Photo:    articleNull.Photo.String,
		Header:   articleNull.Header.String,
		Lead:     articleNull.Lead.String,
		Text:     template.HTML(articleNull.Text.String),
		Date:     articleNull.Date.Time,
		Views:    int(articleNull.Views.Int64),
		Comments: comments,
	}
	return article, nil
}
