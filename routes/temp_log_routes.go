package routes

import (
	"encoding/json"
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
		result := struct {
			Ok bool `json:"ok"`
		}{Ok: true}
		if err := json.NewEncoder(w).Encode(result); err != nil {
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
		var resultError struct {
			Error string `json:"error"`
		}
		err := json.NewDecoder(body).Decode(&dtoAuth)
		if err != nil {
			e = err
			w.WriteHeader(400)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		err = services.VerifyEmail(dtoAuth.Email)
		if err != nil {
			e = err
			w.WriteHeader(400)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		token, err := services.CreateJWT(dtoAuth.Email)
		if err != nil {
			e = err
			w.WriteHeader(400)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		result := struct {
			Token string `json:"token"`
		}{Token: token}
		w.WriteHeader(http.StatusAccepted)
		if err = json.NewEncoder(w).Encode(result); err != nil {
			e = err
		}
		return
	}).Methods("POST")
	return e
}

//postLogs is used for add logs in service
func (ro *Routers) postLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		var resultError struct {
			Error string `json:"error"`
		}
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(401)
			resultError.Error = "unauthorized"
			if err := json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		authorization := r.Header.Get("Authorization")
		payload, err := services.VerifyJWT(authorization)
		if err != nil {
			w.WriteHeader(401)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		body := r.Body
		defer r.Body.Close()
		var dtoLog dto.Log
		transformLog := services.NewTransformLog()
		err = json.NewDecoder(body).Decode(&dtoLog)
		if err != nil {
			e = err
			w.WriteHeader(400)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		transf, err := transformLog.Transform(&dtoLog)
		if err != nil {
			e = err
			w.WriteHeader(400)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}

		ro.Db[payload.Email] = append(ro.Db[payload.Email], transf.Message)
		if len(ro.Db[payload.Email]) > 20 {
			ro.Db[payload.Email] = ro.Db[payload.Email][1:len(ro.Db[payload.Email])]
		}
		w.WriteHeader(201)
		result := struct {
			Status string `json:"status"`
		}{Status: "created"}
		if err = json.NewEncoder(w).Encode(result); err != nil {
			e = err
		}
		return
	}).Methods("POST")
	return e
}

//getLogs return logs of service
func (ro *Routers) getLogs() error {
	var e error = nil
	ro.Mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		var resultError struct {
			Error string `json:"error"`
		}
		if _, ok := r.Header["Authorization"]; !ok {
			w.WriteHeader(401)
			resultError.Error = "unauthorized"
			if err := json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}
		authorization := r.Header.Get("Authorization")
		payload, err := services.VerifyJWT(authorization)
		if err != nil {
			w.WriteHeader(401)
			resultError.Error = err.Error()
			if err = json.NewEncoder(w).Encode(resultError); err != nil {
				e = err
			}
			return
		}

		result := struct {
			Data []string `json:"data"`
		}{Data: ro.Db[payload.Email]}
		w.WriteHeader(http.StatusAccepted)
		if err = json.NewEncoder(w).Encode(result); err != nil {
			e = err
		}
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
