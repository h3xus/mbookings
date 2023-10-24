package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Room struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Meeting struct {
	Title       string `json:"title"`
	RoomID      int    `json:"room_id"`
	UserID      int    `json:"user_id"`
	Status      string `json:"status"`
	BookingDate string `json:"booking_date"`
	Duration    int    `json:"duration"`
}

func getDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":3306)/"+os.Getenv("DB_NAME")+"?parseTime=true")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func addMeetingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, err := getDBConnection()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var meeting Meeting
	err = json.NewDecoder(r.Body).Decode(&meeting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the meeting into the database
	_, err = db.Exec("INSERT INTO meetings (title, room_id, user_id, status, booking_date, duration) VALUES (?, ?, ?, ?, ?, ?)",
		meeting.Title, meeting.RoomID, meeting.UserID, meeting.Status, meeting.BookingDate, meeting.Duration)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/api/users", getUsersHandler)
	http.HandleFunc("/api/rooms", getRoomsHandler)
	http.HandleFunc("/api/meetings", addMeetingHandler)

	port := ":7500"
	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
