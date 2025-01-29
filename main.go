package main

import
(
    "fmt"
    "net/http"
)

func main () {
    mux := &http.ServeMux{}

    mux.HandleFunc("/id/{id}", identify)
    mux.HandleFunc("/", helloworld)

    http.ListenAndServe(":8080", mux)
}

func identify (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintln(w, "Path id:", id)
}

func helloworld (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Hello World!")
}
