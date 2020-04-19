package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
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

	if path.Ext(requestPath) == "" {
		tmpl := template.Must(template.ParseFiles("layouts/file-index.html"))

		if !strings.HasSuffix(requestPath, "/") {
			requestPath += "/"
		}
		objectName, err := url.PathUnescape(strings.Replace(requestPath, "/", "", 1))
		if err != nil {
			log.Fatal(err)
		}
		files := blobHandler.blobService.Files(objectName)

		if requestPath != "/" {
			files = append(files, File{
				Name:        "..",
				IsDirectory: true,
				Path:        path.Dir(strings.TrimSuffix(requestPath, "/")),
			})
		}

		err = tmpl.Execute(w, FilePageData{
			PageTitle: "Media Browser",
			Files:     files,
		})
		if err != nil {
			log.Print(err)
		}
		return
	}
	objectName, err := url.PathUnescape(strings.Replace(requestPath, "/", "", 1))
	if err != nil {
		log.Fatal(err)
	}
	objectName = strings.TrimPrefix(objectName, "/")
	file := blobHandler.blobService.File(objectName)
	w.Header().Add("Location", file.Path)
	w.WriteHeader(302)
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
