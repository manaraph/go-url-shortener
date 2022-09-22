package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	linkList map[string]string
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
	log.Println("Add Link")
	key, ok := r.URL.Query()["link"]
	if ok {
		if !validLink(key[0]) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Could not create shortlink need absolute path link. Ex: /addLink?link=https://github.com/")
			return
		}
		log.Println(key)
		if _, ok := linkList[key[0]]; !ok {
			genString := randStringBytes(10)
			linkList[genString] = key[0]
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusAccepted)

			linkString := fmt.Sprintf("<a href=\"http://localhost:9000/short/%s\">http://localhost:9000/short/%s</a>", genString, genString)
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

// validLink - check that the link we're creating a shortlink for is a absolute URL path
func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	log.Printf("Checking for valid link: %s", link)
	// Check if string matches the regex
	if r.MatchString(link) {
		return true
	}
	return false
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
