package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	articleDeliveryWeb "github.com/mrdniwe/r/internal/article/delivery/web"
	articleRepository "github.com/mrdniwe/r/internal/article/repository/mongo"
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
	v.SetDefault("mongoServer", "localhost")
	v.SetDefault("mongoPort", "27017")
	v.BindEnv("mongoServer", "MONGO_SERVER")
	v.BindEnv("mongoPort", "MONGO_PORT")

	// слушаем события ОС в канал
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)

	// --------
	// подключение к Mongo
	// --------
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + v.GetString("mongoServer") + ":" + v.GetString("mongoPort")))
	if err != nil {
		l.Fatal(err)
	}
	// ждем подключения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		l.Fatal(err)
	}
	// всё ок, можем использовать монго!
	// создаем репозиторий с имеющимся подключением
	articleRepo, err := articleRepository.NewRepository(client, l)
	if err != nil {
		l.Fatal(err)
	}
	// создаем юзкейс с только что созданным репозиторием
	articleUc, err := articleUsecase.NewUsecase(articleRepo, l)
	if err != nil {
		l.Fatal(err)
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
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		client.Disconnect(ctx)
		l.Print("Shutting down")
		os.Exit(0)
	}()
	l.Print("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
