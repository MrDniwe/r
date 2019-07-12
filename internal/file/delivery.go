package filedelivery

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

type FileDelivery struct {
	L *logrus.Logger
	V *viper.Viper
}

func NewDelivery(l *logrus.Logger, r *mux.Router, v *viper.Viper) {
	fd := &FileDelivery{l, v}
	r.HandleFunc("/files/images/{hash}/{name}", fd.Image()).Methods("GET")
	r.HandleFunc("/files/files/{hash}/{name}", fd.File()).Methods("GET")
}

func (fd *FileDelivery) Image() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.Redirect(w, r, fd.V.GetString("s3URIPrefix")+"/images-"+vars["hash"]+"-"+vars["name"], 301)
	}
}

func (fd *FileDelivery) File() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		http.Redirect(w, r, fd.V.GetString("s3URIPrefix")+"/files-"+vars["hash"]+"-"+vars["name"], 301)
	}
}
