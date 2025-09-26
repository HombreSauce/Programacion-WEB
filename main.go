package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", serveIndex)
	port := ":8080"
	fmt.Printf("Servidor escuchando en el puerto http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil) //inicia el servidor
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "./static/index.html")
}
