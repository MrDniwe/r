package usecase

import (
	"github.com/mrdniwe/r/internal/article/repository"
	"github.com/mrdniwe/r/internal/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ArticleUsecase interface {
	SingleArticle(id string) (*models.Article, error)
	LastArticles(amount int) ([]*models.Article, error)
}

type ArticleUC struct {
	Repo repository.ArticleRepository
	L    *logrus.Logger
	V    *viper.Viper
}

func NewUsecase(repo repository.ArticleRepository, l *logrus.Logger, v *viper.Viper) (*ArticleUC, error) {
	return &ArticleUC{repo, l, v}, nil
}

func (u *ArticleUC) SingleArticle(id string) (*models.Article, error) {
	a, err := u.Repo.GetById(id)
	if err != nil {
		return nil, err
	}
	if len(a.Photo) > 0 {
		a.Photo = u.V.GetString("s3URIPrefix") + "/" + a.Photo
	}
	return a, nil
}

func (u *ArticleUC) LastArticles(amount int) ([]*models.Article, error) {
	al, err := u.Repo.GetLastNArticles(amount)
	if err != nil {
		return nil, err
	}
	for _, a := range al {
		if len(a.Photo) > 0 {
			a.Photo = u.V.GetString("s3URIPrefix") + "/" + a.Photo
		}
	}
	return al, nil
}
