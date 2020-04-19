package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type FilePageData struct {
	PageTitle string
	Files     []File
}

type BlobHandler struct {
	blobService BlobService
}

func (blobHandler BlobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.EscapedPath()
	log.Printf("Received request for path %s\n", requestPath)

	tmpl := template.Must(template.ParseFiles("layouts/file-index.html"))
	files := blobHandler.blobService.Files(strings.Replace(requestPath, "/", "", 1))

	if requestPath != "/" {
		files = append(files, File{
			Name:        "..",
			IsDirectory: true,
			Path:        path.Dir(requestPath),
		})
	}

	err := tmpl.Execute(w, FilePageData{
		PageTitle: "Media Browser",
		Files:     files,
	})
	if err != nil {
		log.Print(err)
	}
}

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	handler := BlobHandler{blobService: BlobService{
		context:       ctx,
		storageClient: *client,
		bucketName:    os.Getenv("BUCKET_NAME"),
	}}

	log.Print("Media browser started.")

	http.Handle("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
