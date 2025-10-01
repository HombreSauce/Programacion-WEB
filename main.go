package main

import (
	_ "github.com/lib/pq"
)

func main() {
	//INICIALIZAR SERVIDOR PARA PAGINA WEB
	// http.HandleFunc("/", serveIndex)
	// port := ":8080"
	// fmt.Printf("Servidor escuchando en el puerto http://localhost%s\n", port)
	// err := http.ListenAndServe(port, nil) //inicia el servidor
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// }

	//BASE DE DATOS
	// connStr := "host=localhost port=5432 user=postgres password=postgres dbname=base_turnero sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatalf("failed to connect to DB: %v", err)
	// }
	// defer db.Close()
	// queries := sqlc.New(db)
	// ctx := context.Background()

}

// func serveIndex(w http.ResponseWriter, r *http.Request) {
// 	if r.URL.Path != "/" || r.Method != http.MethodGet {
// 		http.NotFound(w, r)
// 		return
// 	}
// 	w.Header().Set("Content-type", "text/html; charset=utf-8")
// 	http.ServeFile(w, r, "./static/index.html")
// }
