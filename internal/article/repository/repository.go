package repository

import "github.com/mrdniwe/r/internal/models"

type ArticleRepository interface {
	GetById(id string) (*models.Article, error)
	GetLastList(limit, offset int) ([]*models.Article, error)
	PagesCount(int) (int, error)
	UserExists(string) (bool, error)
	NewRecoveryHash(string) (models.RecoveryData, error)
	UserAuth(email, password string) (models.AuthData, error)
	CheckToken(accessToken string) error
	RefreshToken(refreshToken string) (models.AuthData, error)
	LogOutToken(accessToken string) error
	LogOutAllTokens(accessToken string) error
	ChangePassword(email, oldPassword, newPassword string) (models.AuthData, error)
	// RecoveryPassword(uuid, password) error
}
