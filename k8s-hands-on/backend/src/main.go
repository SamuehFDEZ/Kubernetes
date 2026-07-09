package main

import (
    "log"
    "net/http"
)


func ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "hello world"}`))
}

func main() {
    http.HandleFunc("/", ServeHTTP)
    log.Fatal(http.ListenAndServe(":9090", nil))
}