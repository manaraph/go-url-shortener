package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

var (
	linkList map[string]string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()

	linkList = map[string]string{}

	r.HandleFunc("/", showAllLinks)
	r.HandleFunc("/addLink", addLink)
	r.HandleFunc("/short/", getLink)
	r.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port %s", port)
	}

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
	})
	handler := c.Handler(r)
	log.Print("Starting API on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// addLink - Add a link to the linkList and generate a shorter link
// Request eg. http://localhost:8000/addLink?link=https://www.google.com
func addLink(w http.ResponseWriter, r *http.Request) {
	key, ok := r.URL.Query()["link"]
	if ok {
		if _, ok := linkList[key[0]]; !ok {
			genString := fmt.Sprint(rand.Int63n(1000))
			linkList[genString] = key[0]
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusAccepted)
			linkString := fmt.Sprintf("<a href=\"http://localhost:8000/short/%s\">http://localhost:8000/short/%s</a>", genString, genString)
			fmt.Fprintf(w, "Added shortlink\n")
			fmt.Fprintf(w, linkString)
			return
		}
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Already have this link")
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "Failed to add link")
	return
}

// getLink - Find link that matches the shortened link in the linkList
func getLink(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	pathArgs := strings.Split(path, "/")
	log.Printf("Redirected to: %s", linkList[pathArgs[2]])
	http.Redirect(w, r, linkList[pathArgs[2]], http.StatusPermanentRedirect)
	return
}

// Home - Home http request
func showAllLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	log.Println("Get Home")
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	var response string
	for shortLink, link := range linkList {
		response += fmt.Sprintf("Link: <a href=\"http://localhost:8000/short/%s\">http://localhost:8000/short/%s</a> \t\t ShortLink: %s", shortLink, shortLink, link)
	}
	fmt.Fprintf(w, "<h2>Hello and Welcome to the Go URL Shortener!<h2><br>\n")
	fmt.Fprintf(w, response)
	return
}
