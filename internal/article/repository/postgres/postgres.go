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
	return &models.Article{
		Id:     id,
		Header: "Заголовок статьи PG",
		Lead:   "Здесь очень интересная подводка",
		Text:   "Тут сам текст статьи. <b>С тегами!</b><p>И прочая шляпа</p>",
	}, nil
}
