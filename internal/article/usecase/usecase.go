package usecase

import (
	"github.com/mrdniwe/r/internal/article/repository"
	"github.com/mrdniwe/r/internal/models"
	"github.com/sirupsen/logrus"
)

type ArticleUsecase interface {
	SingleArticle(id string) (*models.Article, error)
	LastArticles(amount int) ([]*models.Article, error)
}

type ArticleUC struct {
	Repo repository.ArticleRepository
	L    *logrus.Logger
}

func NewUsecase(repo repository.ArticleRepository, l *logrus.Logger) (*ArticleUC, error) {
	return &ArticleUC{repo, l}, nil
}

func (u *ArticleUC) SingleArticle(id string) (*models.Article, error) {
	return u.Repo.GetById(id)
}

func (u *ArticleUC) LastArticles(amount int) ([]*models.Article, error) {
	return u.Repo.GetLastNArticles(amount)
}
