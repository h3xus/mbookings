package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Room struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Capacity          string `json:"capacity"`
	Feature           string `json:"feature"`
	Availibity_status string `json:"availibity_status"`
	Location          string `json:"location"`
	Created_at        string `json:"created_at"`
}

// type Booking struct {
//  ID            int    `json:"id"`
//  title_meeting string `json:"title_meeting"`
//  room_id       string `json:"room_id"`
//  user_id       string `json:"user_id"`
//  status        string `json:"status"`
//  booking_date  string `json:"booking_date"`
//  duration      string `json:"duration"`
//  created_at    string `json:"created_at"`
// }

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "booker:haslo1@tcp(localhost:3306)/bookings?parseTime=true")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query users from the database
	results, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer results.Close()

	var users []User
	for results.Next() {
		var user User
		err := results.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Respond with the list of users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", "booker:haslo1@tcp(localhost:3306)/bookings?parseTime=true")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query users from the database
	results, err := db.Query("SELECT id, name, location FROM rooms")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer results.Close()

	var rooms []Room
	for results.Next() {
		var room Room
		err := results.Scan(&room.ID, &room.Name, &room.Location)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		rooms = append(rooms, room)
	}

	// Respond with the list of users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func main() {
	http.HandleFunc("/users", getUsersHandler)
	http.HandleFunc("/rooms", getRoomsHandler)
	// http.HandleFunc("/booking", getUsersHandler)

	// Start the server on port 7500
	log.Fatal(http.ListenAndServe(":7500", nil))
}
