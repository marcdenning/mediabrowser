package main

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"time"
)

var (
	ErrFileNotFound = errors.New("specified file could not be found")
)

type File struct {
	Name        string
	IsDirectory bool
	Path        string
}

type BlobStore struct {
	context       context.Context
	storageClient storage.Client
	bucketName    string
}

func (service BlobStore) Files(name string) ([]File, error) {
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
			return nil, err
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
	return files, nil
}

func (service BlobStore) File(name string) (File, error) {
	bucket := service.storageClient.Bucket(service.bucketName)
	object := bucket.Object(name)
	attrs, err := object.Attrs(service.context)
	if err != nil {
		switch err {
		default:
			return File{}, err
		case storage.ErrObjectNotExist:
			return File{}, ErrFileNotFound
		}
	}

	oneDay, err := time.ParseDuration("24h")
	if err != nil {
		return File{}, err
	}

	creds, err := google.FindDefaultCredentials(service.context, storage.ScopeFullControl)
	if err != nil {
		return File{}, err
	}

	conf, err := google.JWTConfigFromJSON(creds.JSON)
	if err != nil {
		return File{}, err
	}

	signedUrl, err := storage.SignedURL(service.bucketName, name, &storage.SignedURLOptions{
		Expires:        time.Now().Add(oneDay),
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Method:         "GET",
		Scheme:         storage.SigningSchemeV4,
	})
	if err != nil {
		return File{}, err
	}
	return File{
		Name: attrs.Name,
		Path: signedUrl,
	}, nil
}
