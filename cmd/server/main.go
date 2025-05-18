package main

import (
	"log"
	"net/http"

	"go-wallet/internal/db"
	"go-wallet/internal/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conn := db.InitPostgres()
	defer conn.Close()

	r := router.Setup(conn)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
