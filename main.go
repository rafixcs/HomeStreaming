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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
