package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	// Template and router init
	r = mux.NewRouter()

}

func main() {
	// Включаем логирование
	logfile, err := os.OpenFile("./log/consolidated.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Не получается открыть log-файл: %v", err)
	}
	defer logfile.Close()
	l = log.New(logfile, "", log.Ldate|log.Ltime)

	// --------
	// подключение к Mongo
	// --------
	var client *mongo.Client
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://mongo:27017"))
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

	// Static
	static := http.FileServer(http.Dir("static"))

	// Middlewares
	// r.Use(mwr["restUri"])

	// Handle and serve
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", static))

	fmt.Println("Server is running on :3000")
	l.Print("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}
