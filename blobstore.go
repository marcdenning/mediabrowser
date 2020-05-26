package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"log"
	"path"
	"time"
)

var (
	ErrFileNotFound = errors.New("specified file could not be found")
)

type File struct {
	Name         string
	IsDirectory  bool
	Path         string
	ContentType  string
	Size         int64
	LastModified time.Time
}

type BlobStore struct {
	context              context.Context
	storageClient        storage.Client
	bucketName           string
	serviceAccountName   string
	privateKeySecretName string
}

func (service BlobStore) Files(name string) ([]File, error) {
	bucket := service.storageClient.Bucket(service.bucketName)
	query := &storage.Query{
		Delimiter: "/",
		Prefix:    name,
	}

	var files []File
	log.Printf("Requesting objects matching %s\n", name)
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
				Name:         path.Base(attrs.Prefix) + "/",
				IsDirectory:  true,
				Path:         "/" + attrs.Prefix,
				ContentType:  attrs.ContentType,
				Size:         attrs.Size,
				LastModified: attrs.Updated,
			})
		} else if attrs.Name != name {
			files = append(files, File{
				Name:         path.Base(attrs.Name),
				Path:         "/" + attrs.Name,
				ContentType:  attrs.ContentType,
				Size:         attrs.Size,
				LastModified: attrs.Updated,
			})
		}
	}
	return files, nil
}

func (service BlobStore) File(name string) (File, error) {
	bucket := service.storageClient.Bucket(service.bucketName)
	object := bucket.Object(name)

	log.Printf("Retrieve attributes for object %s\n", name)
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

	secretClient, err := secretmanager.NewClient(service.context)
	if err != nil {
		return File{}, err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: service.privateKeySecretName,
	}

	result, err := secretClient.AccessSecretVersion(service.context, accessRequest)
	if err != nil {
		return File{}, err
	}

	opts := storage.SignedURLOptions{
		Expires:        time.Now().Add(oneDay),
		GoogleAccessID: service.serviceAccountName,
		PrivateKey:     result.Payload.Data,
		Method:         "GET",
		Scheme:         storage.SigningSchemeV4,
	}

	log.Println("Requesting signed URL for object.")
	signedUrl, err := storage.SignedURL(service.bucketName, name, &opts)
	if err != nil {
		return File{}, err
	}

	return File{
		Name:         attrs.Name,
		Path:         signedUrl,
		ContentType:  attrs.ContentType,
		Size:         attrs.Size,
		LastModified: attrs.Updated,
	}, nil
}
