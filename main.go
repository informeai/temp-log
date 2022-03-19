package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/informeai/temp-log/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
	PORT := os.Getenv("PORT")
	routes := routes.NewRouter()
	if err := routes.Listen(); err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"Origin", "X-Requested-With", "Content-Type", "X-Requested"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD"})
	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(routes.Mux),
		Addr:         fmt.Sprintf(":%s", PORT),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	fmt.Printf("Temp log listening in port: %s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
