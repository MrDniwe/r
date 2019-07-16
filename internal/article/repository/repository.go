package repository

import "github.com/mrdniwe/r/internal/models"

type ArticleRepository interface {
	GetById(id string) (*models.Article, error)
	GetLastList(limit, offset int) ([]*models.Article, error)
}
