package repository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/mrdniwe/r/internal/models"
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

func scanArticle(row scanner) (*models.Article, error) {
	articleNull := articleNullable{}
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
