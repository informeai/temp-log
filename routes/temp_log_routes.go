package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

//Routers is struct for routers.
type Routers struct {
	Mux *mux.Router
}

//NewRouter return instance of Routers
func NewRouter() *Routers {
	return &Routers{Mux: mux.NewRouter()}
}

func (r Routers) health() error {
	var e error = nil
	r.Mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{'ok':'true'}"))
		if err != nil {
			e = err
		}
	})
	return e
}

func (r Routers) pub() error {
	var e error = nil
	r.Mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{'ok':'true'}"))
		if err != nil {
			e = err
		}
	})
	return e
}

//Listen execute listen of routes
func (r Routers) Listen() error {
	var err error
	err = r.health()
	return err
}
