package main

import (
	"fmt"
	"links/config"
	"links/database"
	"links/scan"
	"log"
	"net/http"
	"text/template"
	"time"
)

var (
	tfiles, err = template.ParseFiles("static/index.html")
	index       = template.Must(tfiles, err).Lookup("index")
)

type uploadInjection struct {
	Types []string
}

type indexInjection struct {
	Links []database.Link
}

func serveUpload(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/upload.html")
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, indexInjection{
		Links: database.GetLinks(),
	})
}

func main() {
	log.Println("Starting!")

	handler := http.DefaultServeMux

	static := http.FileServer(http.Dir("./static"))
	handler.Handle("/static/", http.StripPrefix("/static/", static))

	handler.HandleFunc("/", serveIndex)
	handler.HandleFunc("/upload", serveUpload)
	handler.HandleFunc("/upload/new", database.HandleUpload)
	handler.HandleFunc("/scan", scan.HandleScan)

	s := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("127.0.0.1:%d", config.Config.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(s.ListenAndServe())
}
