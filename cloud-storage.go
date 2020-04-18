package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bkt := client.Bucket(os.Getenv("BUCKET_NAME"))

	query := &storage.Query{Prefix: ""}

	var names []string
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
		log.Println(fmt.Sprintf("File: %s\n", attrs.Name))
	}
}
