package mailer

import (
	"github.com/go-gomail/gomail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"strings"
)

type templateSet map[string]*template.Template

type Mailer struct {
	V *viper.Viper
	L *logrus.Logger
	D *gomail.Dialer
	T templateSet
}

type recoveryMailData struct {
	login string
	host  string
	code  string
}

func New(v *viper.Viper, l *logrus.Logger) *Mailer {
	d := gomail.NewDialer(v.GetString("mailServer"), v.GetInt("mailPort"), v.GetString("mailUser"), v.GetString("mailPassword"))
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
	return &Mailer{v, l, d, t}
}

func (m *Mailer) SendRecovery(login string, email string, code string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.V.GetString("mailFrom"))
	msg.SetAddressHeader("To", email, login)
	msg.SetHeader("Subject", "Восстановление пароля")
	var body strings.Builder
	data := recoveryMailData{login, m.V.GetString("code"), code}
	m.T["recovery"].Execute(&body, data)
	msg.SetBody("text/html", body.String())
	if err := m.D.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
