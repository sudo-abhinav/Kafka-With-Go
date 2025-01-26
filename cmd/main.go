package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Sendnotification)
	log.Fatal(http.ListenAndServe(":6900", nil))
}
func Sendnotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
