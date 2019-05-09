package repository

import "github.com/mrdniwe/r/internal/models"

type ArticleRepository interface {
	GetById(id int) (*models.Article, error)
}
