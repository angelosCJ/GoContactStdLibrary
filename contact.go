package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (s *Server) CreateContact(w http.ResponseWriter, r *http.Request) {

	var contact Contact

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = s.db.ExecContext(
		r.Context(),
		`INSERT INTO contacts (id, name, phone, email)
	 VALUES (gen_random_uuid(), $1, $2, $3)`,
		contact.Name,
		contact.Phone,
		contact.Email,
	)

	if err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(contact)

}

func (s *Server) GetContacts(w http.ResponseWriter, r *http.Request) {

	rows, err := s.db.QueryContext(
		r.Context(),
		"SELECT id, name, phone, email FROM contacts",
	)

	if err != nil {
		http.Error(w, "Failed to query contacts", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var contacts []Contact

	for rows.Next() {

		var contact Contact

		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Phone,
			&contact.Email,
		)

		if err != nil {
			http.Error(w, "Failed to read contact", http.StatusInternalServerError)
			return
		}

		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading contacts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(contacts)
}

func (s *Server) GetContactByID(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	var contact Contact

	err := s.db.QueryRowContext(
		r.Context(),
		"SELECT id, name, phone, email FROM contacts WHERE id = $1",
		id,
	).Scan(
		&contact.ID,
		&contact.Name,
		&contact.Phone,
		&contact.Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Contact not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Failed to get contact", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(contact)

}

func (s *Server) UpdateContact(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	var contact Contact

	err := json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := s.db.ExecContext(
		r.Context(),
		`UPDATE contacts
	 SET name = $1, phone = $2, email = $3
	 WHERE id = $4`,
		contact.Name,
		contact.Phone,
		contact.Email,
		id,
	)

	if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contact.ID = id

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(contact)

}

func (s *Server) DeleteContact(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	result, err := s.db.ExecContext(
		r.Context(),
		"DELETE FROM contacts WHERE id = $1",
		id,
	)

	if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
