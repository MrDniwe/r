package filedelivery

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrdniwe/r/internal/server"
)

type FileDelivery struct {
	Srv *server.Server
}

func NewDelivery(r *mux.Router, srv *server.Server) {
	fd := &FileDelivery{srv}
	r.HandleFunc("/files/images/{hash}/{name}", fd.Image()).Methods("GET")
	r.HandleFunc("/files/files/{hash}/{name}", fd.File()).Methods("GET")
}

func (fd *FileDelivery) Image() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.Redirect(w, r, fd.Srv.Conf.GetString("s3URIPrefix")+"/images-"+vars["hash"]+"-"+vars["name"], 301)
	}
}

func (fd *FileDelivery) File() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.Redirect(w, r, fd.Srv.Conf.GetString("s3URIPrefix")+"/files-"+vars["hash"]+"-"+vars["name"], 301)
	}
}
