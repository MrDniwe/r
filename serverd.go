package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	articleDeliveryWeb "github.com/mrdniwe/r/internal/article/delivery/web"
	articleRepository "github.com/mrdniwe/r/internal/article/repository/postgres"
	articleUsecase "github.com/mrdniwe/r/internal/article/usecase"
)

var (
	r *mux.Router
	l *log.Logger
	v *viper.Viper
)

func init() {
	// Подключение конфигурации
	v = viper.New()
	// Настраиваем логгер
	l = log.New()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Template and router init
	r = mux.NewRouter()

}

func main() {
	// определяем конфигурацию
	v.SetDefault("pgHost", "/localhost")
	v.SetDefault("pgPort", "5432")
	v.SetDefault("pgUser", "development")
	v.SetDefault("pgPassword", "development")
	v.SetDefault("pgDbname", "development")
	v.BindEnv("pgHost", "PG_HOST")
	v.BindEnv("pgPort", "PG_PORT")
	v.BindEnv("pgUser", "PG_USER")
	v.BindEnv("pgPassword", "PG_PASSWORD")
	v.BindEnv("pgDbname", "PG_DATABASE")

	// слушаем события ОС в канал
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)

	// --------
	// подключение к БД
	// --------
	connStr := fmt.Sprintf("user=%v dbname=%v sslmode=disable port=%v password=%v host=%v", v.GetString("pgUser"), v.GetString("gUser"), v.GetString("pgDbname"), v.GetString("pgPort"), v.GetString("pgPassword"), v.GetString("pgHost"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
		os.Exit(1)
	}
	// создаем репозиторий с имеющимся подключением
	articleRepo, err := articleRepository.NewRepository(db, l)
	if err != nil {
		l.Fatal(err)
		os.Exit(1)
	}
	// создаем юзкейс с только что созданным репозиторием
	articleUc, err := articleUsecase.NewUsecase(articleRepo, l)
	if err != nil {
		l.Fatal(err)
		os.Exit(1)
	}

	// --------
	// Роуты
	// --------
	//
	// создаем доставку для http
	webRouter := r.PathPrefix("/").Subrouter()
	articleDeliveryWeb.NewDelivery(articleUc, l, webRouter)

	// Handle and serve
	http.Handle("/", r)

	// слушаем события выключения приложения
	go func() {
		sig := <-osChan
		l.Printf("Termination signal --%v-- received", sig)
		db.Close()
		l.Print("Shutting down")
		os.Exit(0)
	}()
	l.Print("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
