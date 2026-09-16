package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

type Server struct {
	db *sql.DB
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
