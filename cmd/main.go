package main

import {
	"fmt"
	"database/sql"
	"log"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	_ "github.com/lib/pq"
}



type Item struct {
	ID int `json:"id"`
	
}

type Order struct {
	ID int `json:"id"`
	Item   string    `json:"item"`
	Amount float64   `json:"amount"`
	Time   time.Time `json:"time"`
}

func main() {
	// connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create router
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/users", createUser(db)).Methods("POST")
	router.HandleFunc("/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser(db)).Methods("DELETE")

	// start server
	log.Fatal(http.listenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type" "application/json")
		next.ServeHTTP(w, r)
	})
}

// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database to get all users
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Create a slice to store the users
		users := []User{}

		// Iterate over the rows and scan each user into the slice
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Username, &user.Password)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		// Convert the users slice to JSON
		jsonData, err := json.Marshal(users)
		if err != nil {
			log.Fatal(err)
		}

		// Set the response header and write the JSON data to the response writer
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

// create a new user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body to get the user data
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response status code and write a success message
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "User created successfully")
	}
}

// get a user by id
func getUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request parameters
		params := mux.Vars(r)
		userID := params["id"]

		// Query the database to get the user by ID
		row := db.QueryRow("SELECT * FROM users WHERE id = $1", userID)

		// Create a User struct to store the retrieved user data
		var user User
		err := row.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the user struct to JSON and write it as the response
		json.NewEncoder(w).Encode(user)
	}
}

// update user by id
func updateUserByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the request parameters
		params := mux.Vars(r)
		userID := params["id"]

		// Parse the request body to get the updated user data
		var updatedUser User
		err := json.NewDecoder(r.Body).Decode(&updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update the user in the database
		_, err = db.Exec("UPDATE users SET username = $1, password = $2 WHERE id = $3", updatedUser.Username, updatedUser.Password, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response status code and write a success message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User updated successfully")
	}
}

// delete user by id

