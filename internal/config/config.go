package config

import (
	"github.com/spf13/viper"
)

func InitialConfig() *viper.Viper {
	// Подключение конфигурации
	v := viper.New()

	// настройки запущенного инстанса
	v.SetDefault("host", "http://localhost:3000")
	v.BindEnv("host", "HOST")

	// настройки БД
	v.SetDefault("pgHost", "localhost")
	v.SetDefault("pgPort", "5434")
	v.SetDefault("pgUser", "development")
	v.SetDefault("pgPassword", "development")
	v.SetDefault("pgDbname", "development")
	v.BindEnv("pgHost", "PG_HOST")
	v.BindEnv("pgPort", "PG_PORT")
	v.BindEnv("pgUser", "PG_USER")
	v.BindEnv("pgPassword", "PG_PASSWORD")
	v.BindEnv("pgDbname", "PG_DATABASE")

	// бакет на S3 где хранятся все файлы контента
	v.SetDefault("s3URIPrefix", "https://r57.s3.eu-central-1.amazonaws.com")
	v.BindEnv("s3URIPrefix", "S3_URI_PREFIX")

	// количество статей на страницу (пагинация)
	v.SetDefault("pageAmount", 10)
	v.BindEnv("pageAmount", "SITE_PAGE_AMOUNT")

	// настройки почты
	v.SetDefault("mailFrom", "robot@region57.com")
	v.SetDefault("mailServer", "smtp.yandex.ru")
	v.SetDefault("mailPort", 465)
	v.SetDefault("mailUser", "robot@region57.com")
	v.SetDefault("mailPassword", "robot14881488")
	v.BindEnv("mailFrom", "MAIL_FROM")
	v.BindEnv("mailServer", "MAIL_SERVER")
	v.BindEnv("mailPort", "MAIL_PORT")
	v.BindEnv("mailUser", "MAIL_USER")
	v.BindEnv("mailPassword", "MAIL_PASSWORD")

	return v
}
