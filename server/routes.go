package server

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	latex "github.com/kevinwylder/pdflatexserver/texwrapper"
)

type TexServer struct {
	index fs.FS
	src   *latex.SourceDirectory
}

func NewLatexServer(src *latex.SourceDirectory, index fs.FS) *TexServer {
	return &TexServer{index, src}
}

func (s *TexServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Println(path)

	filepath, err := s.src.SourcePath(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	info, err := os.Stat(filepath)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if info.IsDir() {
		s.serveDir(w, filepath)
	} else {
		s.servePDF(w, r, filepath)
	}
}

func (s *TexServer) serveDir(w http.ResponseWriter, filepath string) {
	index, err := template.ParseFS(s.index, "template/index.html")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to compile template: %v", err)
		return
	}

	files, err := s.src.ListPath(filepath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = index.Execute(w, files)
	if err != nil {
		fmt.Fprintf(w, "Template error: %v", err)
	}
}

func (s *TexServer) servePDF(w http.ResponseWriter, r *http.Request, filepath string) {
	pdf, err := s.src.PdfCompile(r.Context(), filepath)
	if err != nil {
		switch err.(type) {
		case *latex.CompilerError:
			compilerErr := err.(*latex.CompilerError)
			w.WriteHeader(400)
			io.Copy(w, compilerErr)
		default:
			http.Error(w, err.Error(), 400)
		}
		return
	}

	http.ServeFile(w, r, pdf)
}
