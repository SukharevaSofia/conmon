package main

import (
	"fmt"
	"net/http"
)

func main() {
  initDB()

	http.HandleFunc("/get-table", getTable)
  http.HandleFunc("/ping-upload", postTableUpload)
	http.ListenAndServe(":80", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
