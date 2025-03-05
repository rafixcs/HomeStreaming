package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./views/index.html")
		if err != nil {
			w.Write([]byte("Failed to load index.html"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	})

	fs := http.FileServer(http.Dir("./assets/videos/"))
	http.Handle("/video/", http.StripPrefix("/video/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
