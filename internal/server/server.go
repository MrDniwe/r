package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Server struct {
	Logger *logrus.Logger
	Router *mux.Router
	Conf   *viper.Viper
	Db     *sql.DB
}

func NewServer() *Server {
	conf := config.InitialConfig()

	// настройки логгера
	l := logrus.New()
	switch conf.GetString("logType") {
	case "JSON":
		l.SetFormatter(&logrus.JSONFormatter{})
	default:
		l.SetFormatter(&logrus.TextFormatter{})
	}
	l.SetOutput(os.Stdout)

	// роутер
	r := mux.NewRouter()

	// БД
	connStr := fmt.Sprintf("user=%v dbname=%v sslmode=disable port=%v password=%v host=%v", conf.GetString("pgUser"), conf.GetString("pgDbname"), conf.GetString("pgPort"), conf.GetString("pgPassword"), conf.GetString("pgHost"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		l.Fatal(err)
	}

	return &Server{l, r, conf, db}
}

func (s *Server) ListenAndServe() {
	// слушаем события ОС в канал
	osChan := make(chan os.Signal)
	signal.Notify(osChan, syscall.SIGINT, syscall.SIGTERM)
	// Handle and serve

	http.Handle("/", s.Router)

	// слушаем события выключения приложения
	go func() {
		sig := <-osChan
		s.Logger.Printf("Termination signal --%v-- received", sig)
		s.Db.Close()
		s.Logger.Print("Shutting down")
		os.Exit(0)
	}()
	s.Logger.Print("Server is running on :3000")
	http.ListenAndServe(":"+s.Conf.GetString("appListen"), nil)
}
