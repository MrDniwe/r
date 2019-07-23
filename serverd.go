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
	"github.com/mrdniwe/r/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	articleDeliveryWeb "github.com/mrdniwe/r/internal/article/delivery/web"
	articleRepository "github.com/mrdniwe/r/internal/article/repository/postgres"
	articleUsecase "github.com/mrdniwe/r/internal/article/usecase"
	filesDelivery "github.com/mrdniwe/r/internal/file"
)

var (
	r *mux.Router
	l *log.Logger
	v *viper.Viper
)

func init() {
	// Настраиваем логгер
	l = log.New()
	// l.SetFormatter(&log.JSONFormatter{})
	l.SetFormatter(&log.TextFormatter{})
	l.SetOutput(os.Stdout)
	// Template and router init
	r = mux.NewRouter()

}

func main() {
	// определяем конфигурацию
	v = config.InitialConfig()

	// слушаем события ОС в канал
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)

	// --------
	// подключение к БД
	// --------
	connStr := fmt.Sprintf("user=%v dbname=%v sslmode=disable port=%v password=%v host=%v", v.GetString("pgUser"), v.GetString("pgDbname"), v.GetString("pgPort"), v.GetString("pgPassword"), v.GetString("pgHost"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
	}
	// создаем репозиторий с имеющимся подключением
	articleRepo, err := articleRepository.NewRepository(db, l)
	if err != nil {
		l.Fatal(err)
	}
	// создаем юзкейс с только что созданным репозиторием
	articleUc, err := articleUsecase.NewUsecase(articleRepo, l, v)
	if err != nil {
		l.Fatal(err)
	}

	// --------
	// Роуты
	// --------
	//
	// доставка для пробрасываемых файлов
	filesRouter := r.PathPrefix("/cfs").Subrouter()
	filesDelivery.NewDelivery(l, filesRouter, v)
	// создаем доставку для http
	webRouter := r.PathPrefix("/").Subrouter()
	articleDeliveryWeb.NewDelivery(articleUc, l, webRouter, v)

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
