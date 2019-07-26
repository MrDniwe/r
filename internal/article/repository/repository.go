package repository

import "github.com/mrdniwe/r/internal/models"

type ArticleRepository interface {
	GetById(id string) (*models.Article, error)
	GetLastList(limit, offset int) ([]*models.Article, error)
	PagesCount(int) (int, error)
	UserExists(string) (bool, error)
	NewRecoveryHash(string) (models.RecoveryData, error)
}
