package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"sync"
)

const (
	sourceDir = "/opt/shaleapptest"
)

// implementation of Handler interface that renders a template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join(sourceDir, "templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	cr := newChatRoom()
	fs := http.FileServer(http.Dir(path.Join(sourceDir, "static")))

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.Handle("/chat", cr)

	go cr.run()

	log.Println("Starting web server at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
