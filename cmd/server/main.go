package main

import (
	"log"
	"net/http"

	"wallet-go/internal/db"
	"wallet-go/internal/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conn := db.InitPostgres()
	defer conn.Close()

	r := router.Setup(conn)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
