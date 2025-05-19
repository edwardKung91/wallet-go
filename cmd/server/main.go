package main

import (
	"log"
	"net/http"

	"wallet-go/pkg/db"
	"wallet-go/pkg/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conn := db.InitPostgres()
	defer conn.Close()

	r := router.Setup(conn)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
