package main

import (
	"context"
	"fmt"
    "strings"
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

	mux.HandleFunc("/", helloWorld)
	mux.HandleFunc("/auth", authenticateUser)
	mux.HandleFunc("/animals", getAnimals)
	mux.HandleFunc("/id/{id}", identify)

	http.ListenAndServe(":8080", mux)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func authenticateUser(w http.ResponseWriter, r *http.Request) {
    app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		http.Error(w, "Failed to initialize Firebase app", http.StatusInternalServerError)
		return
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		http.Error(w, "Failed to get Firebase auth client", http.StatusInternalServerError)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	if idToken == authHeader { 
		http.Error(w, "Invalid authorization token format", http.StatusUnauthorized)
		return
	}

	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		http.Error(w, "Invalid or expired ID token", http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"uid":   token.UID,
		"email": token.Claims["email"],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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

