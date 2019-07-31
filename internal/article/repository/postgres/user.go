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

func (a *ArcticleRepo) UserAuth(email, password string) (models.AuthData, error) {
	query := `select email_has_password($1, $2)`
	row := a.Srv.Db.QueryRow(query, email)
	var foundUuid string
	if err := row.Scan(&foundUuid); err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Message {
			case e.NotFoundCode:
				return models.AuthData{}, e.NotFoundErr
			case e.WrongPassword:
				return models.AuthData{}, e.WrongPasswordErr
			default:
				return models.AuthData{}, e.ServerErr
			}
		default:
			return models.AuthData{}, e.ServerErr
		}
	}
	query = `select * from new_tokens($1)`
	row = a.Srv.Db.QueryRow(query, foundUuid)
	var auth models.AuthData
	if err := row.Scan(&auth.AccessToken, &auth.RefreshToken); err != nil {
		return models.AuthData{}, e.ServerErr
	}
	return auth, nil
}
