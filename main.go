package main

import (
	"context"
	"fmt"
	http "net/http"

	// "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var ctx context.Context = context.Background()
var opt option.ClientOption = option.WithCredentialsFile("test-go.json")

func main() {

	mux := &http.ServeMux{}

	mux.HandleFunc("/id/{id}", identify)
	mux.HandleFunc("/", helloworld)
	mux.HandleFunc("/lib", read_library)

	http.ListenAndServe(":8080", mux)
}

func read_library(w http.ResponseWriter, r *http.Request) {
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Fprintln(w, "something went wrong", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Fprintln(w, "something went wrong", err)
	}

	iter := client.Collection("libraries").Doc("library")
	data, err := iter.Get(ctx)
	if err != nil {
		fmt.Fprintln(w, "something went wrong", err)
	}
	fmt.Fprintln(w, "lib?", data.Data())
}

func identify(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintln(w, "Path id:", id)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
