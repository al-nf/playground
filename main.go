package main

import (
	"context"
	"fmt"
    json "encoding/json"
	http "net/http"

	// "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
    "google.golang.org/api/iterator"
)

var ctx context.Context = context.Background()
var opt option.ClientOption = option.WithCredentialsFile("test-go.json")

func main() {

	mux := &http.ServeMux{}

	mux.HandleFunc("/id/{id}", identify)
	mux.HandleFunc("/", helloWorld)
	mux.HandleFunc("/animals", getAnimals)

	http.ListenAndServe(":8080", mux)
}


func getAnimals(w http.ResponseWriter, r *http.Request) {
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Fprintln(w, "something went wrong", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Fprintln(w, "something went wrong", err)
	}

	iter := client.Collection("animals").Documents(ctx)

    type Animal struct {
		ID   string `json:"id"`
		Extinct bool `json:"extinct"`
	}

    var animals []Animal

    for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading document: %v", err), http.StatusInternalServerError)
			return
		}

		extinct, ok := doc.Data()["extinct"].(bool)
		if !ok {
			extinct = false
		}

		animals = append(animals, Animal{
			ID:      doc.Ref.ID,
			Extinct: extinct,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(animals); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func identify(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintln(w, "Path id:", id)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
