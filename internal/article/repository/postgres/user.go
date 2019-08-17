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
				a.Srv.Logger.WithFields(logrus.Fields{
					"type": e.PostgresError,
				}).Error(err)
				return models.RecoveryData{}, e.ServerErr

			}
		default:
			if err == sql.ErrNoRows {
				return models.RecoveryData{}, e.NotFoundErr
			}
			a.Srv.Logger.WithFields(logrus.Fields{
				"type": e.PostgresError,
			}).Error(err)
			return models.RecoveryData{}, e.ServerErr
		}
	}
	return models.RecoveryData{login, email, hash}, nil
}

func (a *ArcticleRepo) UserAuth(email, password string) (models.AuthData, error) {
	query := `select email_has_password($1, $2)`
	row := a.Srv.Db.QueryRow(query, email, password)
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
				a.Srv.Logger.WithFields(logrus.Fields{
					"type": e.PostgresError,
				}).Error(err)
				return models.AuthData{}, e.ServerErr
			}
		default:
			a.Srv.Logger.WithFields(logrus.Fields{
				"type": e.PostgresError,
			}).Error(err)
			return models.AuthData{}, e.ServerErr
		}
	}
	query = `select * from new_tokens($1)`
	row = a.Srv.Db.QueryRow(query, foundUuid)
	var auth models.AuthData
	if err := row.Scan(&auth.AccessToken, &auth.RefreshToken); err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return models.AuthData{}, e.ServerErr
	}
	return auth, nil
}

func (a *ArcticleRepo) CheckToken(accessToken string) error {
	query := `select is_valid_access_token($1)`
	var isValid bool
	row := a.Srv.Db.QueryRow(query, accessToken)
	if err := row.Scan(&isValid); err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return e.ServerErr
	}
	if !isValid {
		return e.InvalidTokenErr
	}
	return nil
}

func (a *ArcticleRepo) RefreshToken(refreshToken string) (models.AuthData, error) {
	query := `select * from refresh_tokens($1)`
	var auth models.AuthData
	row := a.Srv.Db.QueryRow(query, refreshToken)
	if err := row.Scan(&auth.AccessToken, &auth.RefreshToken); err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Message {
			case e.TokenNotFound:
				return models.AuthData{}, e.InvalidTokenErr
			default:
				a.Srv.Logger.WithFields(logrus.Fields{
					"type": e.PostgresError,
				}).Error(err)
				return models.AuthData{}, e.ServerErr
			}
		default:
			a.Srv.Logger.WithFields(logrus.Fields{
				"type": e.PostgresError,
			}).Error(err)
			return models.AuthData{}, e.ServerErr
		}
	}
	return auth, nil
}

func (a *ArcticleRepo) LogOutToken(accessToken string) error {
	query := `delete from tokens where access_token =$1`
	_, err := a.Srv.Db.Exec(query, accessToken)
	if err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return e.ServerErr
	}
	return nil
}

func (a *ArcticleRepo) LogOutAllTokens(accessToken string) error {
	query := `delete from tokens where user_uuid = (select user_uuid from tokens where access_token = $1)`
	_, err := a.Srv.Db.Exec(query, accessToken)
	if err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return e.ServerErr
	}
	return nil
}

func (a *ArcticleRepo) LogOutAllTokensByEmail(email string) error {
	query := `delete from tokens where user_uuid = (select uuid from users where email = $1)`
	_, err := a.Srv.Db.Exec(query, email)
	if err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return e.ServerErr
	}
	return nil
}

func (a *ArcticleRepo) ChangePassword(email, oldPassword, newPassword string) (models.AuthData, error) {
	// проверяем правильность старого пароля
	query := `select email_has_password($1, $2)`
	row := a.Srv.Db.QueryRow(query, email, oldPassword)
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
				a.Srv.Logger.WithFields(logrus.Fields{
					"type": e.PostgresError,
				}).Error(err)
				return models.AuthData{}, e.ServerErr
			}
		default:
			a.Srv.Logger.WithFields(logrus.Fields{
				"type": e.PostgresError,
			}).Error(err)
			return models.AuthData{}, e.ServerErr
		}
	}
	// меняем старый пароль на новый
	query = `update users set password =$1 where email = $2`
	_, err := a.Srv.Db.Exec(query, newPassword, email)
	if err != nil {
		a.Srv.Logger.WithFields(logrus.Fields{
			"type": e.PostgresError,
		}).Error(err)
		return models.AuthData{}, e.ServerErr
	}
	// отзываетм все токены
	err = a.LogOutAllTokensByEmail(email)
	if err != nil {
		return models.AuthData{}, err
	}
	// авторизуемся с новым паролем и возвращаем новые токены
	return a.UserAuth(email, newPassword)
}
