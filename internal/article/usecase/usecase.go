package usecase

import (
	"github.com/mrdniwe/r/internal/article/repository"
	"github.com/mrdniwe/r/internal/models"
	"github.com/mrdniwe/r/internal/server"
	e "github.com/mrdniwe/r/pkg/errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
)

type ArticleUsecase interface {
	SingleArticle(id string) (*models.Article, error)
	LastArticles(amount, offset int) ([]*models.Article, error)
	TotalPagesCount() (int, error)
	CheckEmailExists(string) (bool, error)
	NewRecoveryHash(string) (models.RecoveryData, error)
	UserAuth(email, password string) (models.AuthData, error)
}

type ArticleUC struct {
	Repo repository.ArticleRepository
	Srv  *server.Server
}

func NewUsecase(repo repository.ArticleRepository, srv *server.Server) (*ArticleUC, error) {
	return &ArticleUC{repo, srv}, nil
}

func (u *ArticleUC) SingleArticle(id string) (*models.Article, error) {
	// валидируем uuid
	uu, err := uuid.ParseHex(id)
	if err != nil {
		u.Srv.Logger.WithFields(logrus.Fields{
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
		a.Photo = u.Srv.Conf.GetString("s3URIPrefix") + "/" + a.Photo
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
			a.Photo = u.Srv.Conf.GetString("s3URIPrefix") + "/" + a.Photo
		}
	}
	return al, nil
}

func (u *ArticleUC) TotalPagesCount() (int, error) {
	total, err := u.Repo.PagesCount(u.Srv.Conf.GetInt("pageAmount"))
	if err != nil {
		return 0, err
	}
	return total, nil
}
