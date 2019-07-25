package mailer

import (
	"html/template"
	"strings"

	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type templateSet map[string]*template.Template

type Mailer struct {
	Logger    *logrus.Logger
	Conf      *viper.Viper
	Dialer    *gomail.Dialer
	Templates templateSet
}

type recoveryMailData struct {
	login string
	host  string
	code  string
}

func New(l *logrus.Logger, c *viper.Viper) *Mailer {
	d := gomail.NewPlainDialer(c.GetString("mailServer"), c.GetInt("mailPort"), c.GetString("mailUser"), c.GetString("mailPassword"))
	t := templateSet{}
	t["recovery"] = template.New(`
	<h1>Восстановление пароля</h1>
	<p>
		Вы запросили восстановление пароля на сайте region57 для логина {{.login}}
	</p>
	<p>
		Для восстановления пароля пройдите 
		<a href="{{.host}}/recovery/{{.code}}" target="_blank">по ссылке</a>
	</p>
	<p>
		Если вы не запрашивали восстановления, вероятно, это действия злоумышленников.
		В таком случае, переходить по ссылке не нужно.
	</p>
	`)
	return &Mailer{l, c, d, t}
}

func (m *Mailer) SendRecovery(login string, email string, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Conf.GetString("mailFrom"))
	msg.SetAddressHeader("To", email, login)
	msg.SetHeader("Subject", "Восстановление пароля")
	var body strings.Builder
	data := recoveryMailData{login, m.Conf.GetString("code"), code}
	m.Templates["recovery"].Execute(&body, data)
	msg.SetBody("text/html", body.String())
	if err := m.Dialer.DialAndSend(msg); err != nil {
		m.Logger.WithFields(logrus.Fields{
			"type": e.MailError,
		}).Error(err)
		return e.ServerErr
	}
	return nil
}
