package main

import (
	"cloud.google.com/go/storage"
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type File struct {
	Name        string
	IsDirectory bool
	Path        string
}

type BlobService struct {
	context       context.Context
	storageClient storage.Client
	bucketName    string
}

func (service BlobService) Files(name string) []File {
	bucket := service.storageClient.Bucket(service.bucketName)
	query := &storage.Query{
		Delimiter: "/",
		Prefix:    name,
	}

	var files []File
	it := bucket.Objects(service.context, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if attrs.Prefix != "" {
			files = append(files, File{
				Name:        attrs.Prefix,
				IsDirectory: true,
				Path:        "/" + attrs.Prefix,
			})
		} else if attrs.Name != name {
			files = append(files, File{
				Name: attrs.Name,
				Path: "/" + attrs.Name,
			})
		}
	}
	return files
}

func (service BlobService) File(name string) File {
	bucket := service.storageClient.Bucket(service.bucketName)
	object := bucket.Object(name)
	attrs, err := object.Attrs(service.context)
	if err != nil {
		log.Fatal(err)
	}

	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		log.Fatal(err)
	}

	jsonKey, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		log.Fatal(err)
	}

	signedUrl, err := storage.SignedURL(service.bucketName, name, &storage.SignedURLOptions{
		Expires:        time.Now().Add(oneDay),
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Method:         "GET",
		Scheme:         storage.SigningSchemeV4,
	})
	if err != nil {
		log.Fatal(err)
	}
	return File{
		Name: attrs.Name,
		Path: signedUrl,
	}
}
