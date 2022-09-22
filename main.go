package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var NewUrl string = "https://github.com"

func main() {
	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	http.Handle("/", http.RedirectHandler(NewUrl, http.StatusMovedPermanently))
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		log.Printf("Defaulting to port %s", port)
	}

	// c := cors.New(cors.Options{
	// 	AllowCredentials: true,
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
	// 	AllowedHeaders:   []string{"*"},
	// })
	// handler := c.Handler(r)
	log.Print("Starting API on 0.0.0.0:" + port)
	// log.Fatal(http.ListenAndServe(":"+port, handler))
	log.Fatal(http.ListenAndServe(":9000", nil))

}
