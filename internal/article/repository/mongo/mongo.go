package repository

import (
	"context"
	"github.com/mrdniwe/r/internal/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRepository(client *mongo.Client, l *logrus.Logger) (*ArcticleRepo, error) {
	// пингуем на всякий случай
	err := client.Ping(context.Background(), nil)
	if err != nil {
		l.Fatal(err)
	}
	return &ArcticleRepo{client, l}, nil
}

type ArcticleRepo struct {
	db *mongo.Client
	L  *logrus.Logger
}

func (a *ArcticleRepo) GetById(id int) (*models.Article, error) {
	return &models.Article{
		Id:     id,
		Header: "Заголовок статьи",
		Pre:    "Здесь очень интересная подводка",
		Text:   "Тут сам текст статьи. <b>С тегами!</b><p>И прочая шляпа</p>",
	}, nil
}
