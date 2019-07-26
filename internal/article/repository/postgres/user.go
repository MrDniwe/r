package repository

import (
	"github.com/mrdniwe/r/internal/models"
	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (a *ArcticleRepo) UserExists(email string) (bool, error) {
	query := `
		select count(*) as count
		from users
		where lower($1) = lower(email)
		`
	row := a.Srv.Db.QueryRow(query, email)
	var found int
	if err := row.Scan(&found); err != nil {
		nerr := errors.Wrap(err, "Error while searching for user by email")
		if err, ok := nerr.(e.StackTracer); ok {
			a.Srv.Logger.WithFields(logrus.Fields{
				"type":  e.PostgresError,
				"stack": err.StackTrace()[0],
			}).Error(err)
		}
	}
	if found < 1 {
		return false, nil
	}
	return true, nil
}

func (a *ArcticleRepo) NewRecoveryHash(email string) (models.RecoveryData, error) {
	return models.RecoveryData{"Lasso", "ceperagrey@gmail.com", "qazwsx"}, nil
}
