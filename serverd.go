package main

import (
	_ "github.com/lib/pq"

	articleDeliveryWeb "github.com/mrdniwe/r/internal/article/delivery/web"
	articleRepository "github.com/mrdniwe/r/internal/article/repository/postgres"
	articleUsecase "github.com/mrdniwe/r/internal/article/usecase"
	filesDelivery "github.com/mrdniwe/r/internal/file"
	"github.com/mrdniwe/r/internal/server"
)

var (
	srv *server.Server
)

func init() {
	srv = server.NewServer()
}

func main() {

	// создаем репозиторий с имеющимся подключением
	articleRepo, err := articleRepository.NewRepository(srv)
	if err != nil {
		srv.Logger.Fatal(err)
	}
	// создаем юзкейс с только что созданным репозиторием
	articleUc, err := articleUsecase.NewUsecase(articleRepo, srv)
	if err != nil {
		srv.Logger.Fatal(err)
	}

	// --------
	// Роуты
	// --------
	//
	// доставка для пробрасываемых файлов
	filesRouter := srv.Router.PathPrefix("/cfs").Subrouter()
	filesDelivery.NewDelivery(filesRouter, srv)
	// создаем доставку для http
	webRouter := srv.Router.PathPrefix("/").Subrouter()
	articleDeliveryWeb.NewDelivery(articleUc, webRouter, srv)

	srv.ListenAndServe()
}
