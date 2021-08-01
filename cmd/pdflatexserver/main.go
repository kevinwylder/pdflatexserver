package main

import (
	"log"
	"net/http"

	"github.com/kevinwylder/pdflatexserver/server"
	latex "github.com/kevinwylder/pdflatexserver/texwrapper"
)

func main() {
	src, err := latex.NewSourceDirectory("/data")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", server.NewLatexServer(src, "/srv/http/index.html"))
	log.Println("Starting server")
	http.ListenAndServe(":80", nil)
}
