package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./views/home/index.html")
		if err != nil {
			w.Write([]byte("Failed to load home"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/videos/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./views/videos/index.html")
		if err != nil {
			w.Write([]byte("Failed to load videos page"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	})

	fs := http.FileServer(http.Dir("./assets/videos/"))
	http.Handle("/video/", http.StripPrefix("/video/", fs))

	// Serve static files (JS, CSS, images, etc.)
	fs = http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
