package usecase

import (
	"github.com/mrdniwe/r/internal/article/repository"
	"github.com/mrdniwe/r/internal/models"
	"log"
)

type ArticleUsecase interface {
	SingleArticle(id int) (*models.Article, error)
	LastArticles(amount int) (*[]models.Article, error)
}

type ArticleUC struct {
	Repo repository.ArticleRepository
	L    *log.Logger
}

func NewUsecase(repo repository.ArticleRepository, l *log.Logger) (*ArticleUC, error) {
	return &ArticleUC{repo, l}, nil
}

func (u *ArticleUC) SingleArticle(id int) (*models.Article, error) {
	return u.Repo.GetById(id)
}
