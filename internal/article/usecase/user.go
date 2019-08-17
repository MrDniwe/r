package usecase

import (
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
)

func (u *ArticleUC) CheckEmailExists(email string) (bool, error) {
	return u.Repo.UserExists(email)
}

func (u *ArticleUC) NewRecoveryHash(email string) (models.RecoveryData, error) {
	return u.Repo.NewRecoveryHash(email)
}

func (u *ArticleUC) UserAuth(email, password string) (models.AuthData, error) {
	return u.Repo.UserAuth(email, password)
}

func (u *ArticleUC) UserLogout(accessToken string) error {
	return u.Repo.LogOutToken(accessToken)
}

func (u *ArticleUC) CheckAccessToken(accessToken string) error {
	if accessToken == "" {
		return e.InvalidTokenErr
	}
	return u.Repo.CheckToken(accessToken)
}

func (u *ArticleUC) RefreshToken(refreshToken string) (models.AuthData, error) {
	return u.Repo.RefreshToken(refreshToken)
}
