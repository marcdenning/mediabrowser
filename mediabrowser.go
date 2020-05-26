package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/c2h5oh/datasize"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type FilePageData struct {
	PageTitle string
	Files     []File
}

func formatTimeForHTML(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

func formatBytesForHTML(b int64) string {
	return datasize.ByteSize(b).HR()
}

func authenticateHandler(handler func(w http.ResponseWriter, r *http.Request), username, password string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		providedUser, providedPass, ok := r.BasicAuth()
		if !ok {
			log.Println("Could not find or parse Authorization header.")
			w.Header().Add("WWW-Authenticate", "Basic realm=\"mediabrowser\"")
			http.Error(w, "Could not find or parse Authorization header.", http.StatusUnauthorized)
			return
		}
		if username != providedUser || password != providedPass {
			log.Println("Invalid credentials provided.")
			w.Header().Add("WWW-Authenticate", "Basic realm=\"mediabrowser\"")
			http.Error(w, "Invalid credentials provided.", http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func serveBlobs(service BlobStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.EscapedPath()
		log.Printf("Received request for path %s\n", requestPath)

		if path.Ext(requestPath) == "" {
			tmpl, err := template.New("layouts/file-index.html").Funcs(template.FuncMap{
				"formatTime":  formatTimeForHTML,
				"formatBytes": formatBytesForHTML,
			}).ParseFiles("layouts/file-index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !strings.HasSuffix(requestPath, "/") {
				requestPath += "/"
			}
			objectName, err := url.PathUnescape(strings.Replace(requestPath, "/", "", 1))
			if err != nil {
				log.Println("Could not parse request path.", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			files, err := service.Files(objectName)

			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			switch r.Method {
			case http.MethodHead:
				w.Header().Add("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
			case http.MethodGet:
				if requestPath != "/" {
					files = append([]File{
						{
							Name:        "..",
							IsDirectory: true,
							Path:        path.Dir(strings.TrimSuffix(requestPath, "/")),
						},
					}, files...)
				}

				err = tmpl.Execute(w, FilePageData{
					PageTitle: fmt.Sprintf("Index of %s", objectName),
					Files:     files,
				})
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			return
		}

		objectName, err := url.PathUnescape(strings.Replace(requestPath, "/", "", 1))
		if err != nil {
			log.Println("Could not parse request path.", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		objectName = strings.TrimPrefix(objectName, "/")
		file, err := service.File(objectName)

		if err != nil {
			log.Println("Could not retrieve file info.", err)
			switch err {
			case ErrFileNotFound:
				http.NotFound(w, r)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		switch r.Method {
		case http.MethodHead:
			w.Header().Add("Content-Type", file.ContentType)
			w.Header().Add("Content-Length", fmt.Sprintf("%d", file.Size))
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			http.Redirect(w, r, file.Path, http.StatusFound)
		}
	}
}

func serveStaticFile(name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	blobService := BlobStore{
		context:              ctx,
		storageClient:        *client,
		bucketName:           os.Getenv("BUCKET_NAME"),
		serviceAccountName:   os.Getenv("SERVICE_ACCOUNT_NAME"),
		privateKeySecretName: os.Getenv("PK_SECRET_NAME"),
	}

	log.Print("Media browser started.")

	var handler http.HandlerFunc
	username, ok := os.LookupEnv("WEB_USERNAME")
	password, ok := os.LookupEnv("WEB_PASSWORD")

	if !ok {
		log.Println("Variables WEB_USERNAME and WEB_PASSWORD not set. Will not authenticate requests.")
		handler = serveBlobs(blobService)
	} else {
		handler = authenticateHandler(serveBlobs(blobService), username, password)
	}

	http.HandleFunc("/favicon.ico", serveStaticFile("favicon.ico"))
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
