package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	article_delivery_web "github.com/mrdniwe/r/internal/article/delivery/web"
	article_repository "github.com/mrdniwe/r/internal/article/repository/mongo"
	article_usecase "github.com/mrdniwe/r/internal/article/usecase"
)

// global app vars
var (
	r *mux.Router
	l *log.Logger
)

func init() {
	// Настраиваем логгер
	l = log.New()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Template and router init
	r = mux.NewRouter()

}

func main() {

	// --------
	// подключение к Mongo
	// --------
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		l.Fatal(err)
	}
	// ждем подключения
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		l.Fatal(err)
	}
	// всё ок, можем использовать монго!
	// создаем репозиторий с имеющимся подключением
	article_repo, err := article_repository.NewRepository(client, l)
	if err != nil {
		l.Fatal(err)
	}
	// создаем юзкейс с только что созданным репозиторием
	article_uc, err := article_usecase.NewUsecase(article_repo, l)
	if err != nil {
		l.Fatal(err)
	}

	// --------
	// Роуты
	// --------
	//
	// создаем доставку для http
	web_router := r.PathPrefix("/").Subrouter()
	article_delivery_web.NewDelivery(article_uc, l, web_router)

	// создаем доставку для api
	// TODO

	// Middlewares
	// r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)

	fmt.Println("Server is running on :3000")
	l.Print("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
