package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/informeai/temp-log/database"
	"github.com/informeai/temp-log/dto"
	"github.com/informeai/temp-log/services"
)

//Routers is struct for routers.
type Routers struct {
	Mux *mux.Router
	Db  database.Mock
}

//NewRouter return instance of Routers
func NewRouter() *Routers {
	return &Routers{Mux: mux.NewRouter(), Db: database.Mock{}}
}

//health is return status of service
func (ro *Routers) health() error {
	var e error = nil
	ro.Mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{'ok':'true'}"))
		if err != nil {
			e = err
		}
	})
	return e
}

//postLogs is used for add logs in service
func (ro *Routers) postLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer r.Body.Close()
		var dtoLog dto.Log
		transformLog := services.NewTransformLog()
		err := json.NewDecoder(body).Decode(&dtoLog)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		transf, err := transformLog.Transform(&dtoLog)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		if len(ro.Db) > 20 {
			ro.Db = ro.Db[1:len(ro.Db)]
		}
		ro.Db = append(ro.Db, transf)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(transf.Message))
		return
	}).Methods("POST")
	return e
}

//getLogs return logs of service
func (ro *Routers) getLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		mockBytes, err := json.Marshal(ro.Db)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(mockBytes)
		return
	}).Methods("GET")
	return e
}

//Listen execute listen of routes
func (ro *Routers) Listen() error {
	var err error
	err = ro.health()
	err = ro.getLogs()
	err = ro.postLogs()
	return err
}
