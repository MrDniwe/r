package usecase

import "github.com/mrdniwe/r/internal/models"

func (u *ArticleUC) CheckEmailExists(email string) (bool, error) {
	return u.Repo.UserExists(email)
}

func (u *ArticleUC) NewRecoveryHash(email string) (models.RecoveryData, error) {
	return u.Repo.NewRecoveryHash(email)
}

func (u *ArticleUC) UserAuth(email, password string) (models.AuthData, error) {
	return u.Repo.UserAuth(email, password)
}
