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

//auth is used for return token of access service
func (ro *Routers) auth() error {
	var e error = nil
	ro.Mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer r.Body.Close()
		var dtoAuth dto.Auth
		err := json.NewDecoder(body).Decode(&dtoAuth)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		err = services.VerifyEmail(dtoAuth.Email)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		token, err := services.CreateJWT(dtoAuth.Email)
		if err != nil {
			e = err
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf("{'token': %s}", token)))
		return
	}).Methods("POST")
	return e
}

//postLogs is used for add logs in service
func (ro *Routers) postLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(401)
			w.Write([]byte(fmt.Sprint("{error: unauthorized}")))
			return
		}
		authorization := r.Header.Get("Authorization")
		payload, err := services.VerifyJWT(authorization)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		fmt.Printf("payload: %+s\n", payload)
		body := r.Body
		defer r.Body.Close()
		var dtoLog dto.Log
		transformLog := services.NewTransformLog()
		err = json.NewDecoder(body).Decode(&dtoLog)
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

		ro.Db[payload.Email] = append(ro.Db[payload.Email], transf.Message)
		w.WriteHeader(201)
		w.Write([]byte("{status: created}"))
		return
	}).Methods("POST")
	return e
}

//getLogs return logs of service
func (ro *Routers) getLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(401)
			w.Write([]byte(fmt.Sprint("{error: unauthorized}")))
			return
		}
		authorization := r.Header.Get("Authorization")
		payload, err := services.VerifyJWT(authorization)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(fmt.Sprintf("{error: %s}", err.Error())))
			return
		}
		fmt.Printf("payload: %+s\n", payload)

		mockBytes, err := json.Marshal(ro.Db[payload.Email])
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
	err = ro.auth()
	err = ro.health()
	err = ro.getLogs()
	err = ro.postLogs()
	return err
}
