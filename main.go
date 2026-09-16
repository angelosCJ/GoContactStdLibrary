package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type Server struct {
	db *sql.DB
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	server := &Server{
		db: db,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /contacts", server.GetContacts)
	mux.HandleFunc("GET /contacts/{id}", server.GetContactByID)
	mux.HandleFunc("POST /contacts", server.CreateContact)
	mux.HandleFunc("PUT /contacts/{id}", server.UpdateContact)
	mux.HandleFunc("DELETE /contacts/{id}", server.DeleteContact)

	log.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
