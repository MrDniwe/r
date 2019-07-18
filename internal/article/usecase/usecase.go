package usecase

import (
	"github.com/mrdniwe/r/internal/article/repository"
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ArticleUsecase interface {
	SingleArticle(id string) (*models.Article, error)
	LastArticles(amount, offset int) ([]*models.Article, error)
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
	// валидируем uuid
	uu, err := uuid.ParseHex(id)
	if err != nil {
		u.L.WithFields(logrus.Fields{
			"type":  e.ValidationError,
			"in":    "UUID",
			"given": id,
		}).Info(err)
		return nil, e.BadRequestErr
	}
	a, err := u.Repo.GetById(uu.String())
	if err != nil {
		return nil, err
	}
	if len(a.Photo) > 0 {
		a.Photo = u.V.GetString("s3URIPrefix") + "/" + a.Photo
	}
	return a, nil
}

func (u *ArticleUC) LastArticles(amount, offset int) ([]*models.Article, error) {
	al, err := u.Repo.GetLastList(amount, offset)
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
