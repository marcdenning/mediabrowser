package main

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
	"log"
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
		files = append(files, File{
			Name: attrs.Name,
			Path: "/" + attrs.Name,
		})
	}
	return files
}
