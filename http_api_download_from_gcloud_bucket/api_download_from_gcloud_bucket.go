package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
)

func main() {
	router := httprouter.New()
	router.GET("/download", Download)

	log.Fatal(http.ListenAndServe(":8000", router))
}

func Download(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/path/to/credential.json"))
	if err != nil {
		fmt.Println("client errorrr", err)
	}
	bucket := client.Bucket("bucket-name")

	rc, err := bucket.Object("/your/object/to-download.sql").NewReader(ctx)
	if err != nil {
		fmt.Println("rc errorrr", err)
	}
	defer rc.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=haha.sql")
	_, err = io.Copy(w, rc)
	if err != nil {
		fmt.Println("io copy errr", err)
	}
}
