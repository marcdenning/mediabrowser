package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type File struct {
	Name        string
	IsDirectory bool
}

type FilePageData struct {
	PageTitle string
	Files     []File
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Received request")
	tmpl := template.Must(template.ParseFiles("layouts/file-index.html"))

	err := tmpl.Execute(w, FilePageData{
		PageTitle: "Media Browser",
		Files: []File{
			{Name: "sample.txt", IsDirectory: false},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

func main() {
	log.Print("Media browser started.")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
