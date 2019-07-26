package repository

import (
	"database/sql"

	"github.com/lib/pq"
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
	query := `select create_recovery_hash($1) as hash, login from users u where lower($1) = lower(u.email);`
	row := a.Srv.Db.QueryRow(query, email)
	var hash, login string
	if err := row.Scan(&hash, &login); err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Message {
			case e.ToSoonCode:
				return models.RecoveryData{}, e.DelayErr
			case e.NotFoundCode:
				return models.RecoveryData{}, e.NotFoundErr
			default:
				return models.RecoveryData{}, e.ServerErr

			}
		default:
			if err == sql.ErrNoRows {
				return models.RecoveryData{}, e.NotFoundErr
			}
			return models.RecoveryData{}, e.ServerErr
		}
	}
	return models.RecoveryData{login, email, hash}, nil
}
