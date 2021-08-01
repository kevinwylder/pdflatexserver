package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/kevinwylder/pdflatexserver/server"
	latex "github.com/kevinwylder/pdflatexserver/texwrapper"
)

//go:embed template/index.html
var embeddedTemplate embed.FS

var (
	templateDir = flag.String("template", "", "path to template dir")
	port        = flag.String("port", ":80", "port")
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("USAGE: pdflatexserver <flags> <directory to serve>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	src, err := latex.NewSourceDirectory(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	var data fs.FS
	if *templateDir == "" {
		data = embeddedTemplate
	} else {
		data = os.DirFS(*templateDir)
	}
	http.Handle("/", server.NewLatexServer(src, data))
	log.Printf("Starting server on %s\n", *port)
	http.ListenAndServe(*port, nil)
}
