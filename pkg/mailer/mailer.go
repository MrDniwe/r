package mailer

import (
	"html/template"
	"strings"

	e "github.com/mrdniwe/r/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	Logger    *logrus.Logger
	Conf      *viper.Viper
	Dialer    *gomail.Dialer
	Templates *template.Template
}

type recoveryMailData struct {
	Login string
	Host  string
	Code  string
}

func New(l *logrus.Logger, c *viper.Viper) *Mailer {
	d := gomail.NewPlainDialer(c.GetString("mailServer"), c.GetInt("mailPort"), c.GetString("mailUser"), c.GetString("mailPassword"))
	t, _ := template.New("recovery").Parse(`
	<h1>Восстановление пароля</h1>
	<p>
		Вы запросили восстановление пароля на сайте Region57 для логина <b>{{.Login}}</b>
	</p>
	<p>
		Для восстановления пароля пройдите по ссылке
		<a href="{{.Host}}/recovery-submit?code={{.Code}}" target="_blank">{{.Host}}/recovery-submit</a>, введите в поле этот код <i>{{.Code}}</i>
		и запросите смену пароля.
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
	data := recoveryMailData{login, m.Conf.GetString("host"), code}
	m.Templates.ExecuteTemplate(&body, "recovery", data)
	msg.SetBody("text/html", body.String())
	if err := m.Dialer.DialAndSend(msg); err != nil {
		m.Logger.WithFields(logrus.Fields{
			"type": e.MailError,
		}).Error(err)
		return e.ServerErr
	}
	return nil
}
