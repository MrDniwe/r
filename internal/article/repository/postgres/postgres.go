package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
	"github.com/sirupsen/logrus"
	"html/template"
)

func NewRepository(client *sql.DB, l *logrus.Logger) (*ArcticleRepo, error) {
	return &ArcticleRepo{client, l}, nil
}

type ArcticleRepo struct {
	db *sql.DB
	L  *logrus.Logger
}

type ArticleNullable struct {
	Id     string
	Header sql.NullString
	Lead   sql.NullString
	Text   sql.NullString
	Date   pq.NullTime
	Photo  sql.NullString
	Views  sql.NullInt64
}

func (a *ArcticleRepo) GetById(id string) (*models.Article, error) {
	query := `select 
	  uuid, title, lead, body, active_from, views, image
	from articles
	  where is_visible=true
	  and uuid=$1`
	row := a.db.QueryRow(query, id)
	articleNull := ArticleNullable{}
	if err := row.Scan(&articleNull.Id, &articleNull.Header, &articleNull.Lead, &articleNull.Text, &articleNull.Date, &articleNull.Views, &articleNull.Photo); err != nil {
		return nil, err
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
