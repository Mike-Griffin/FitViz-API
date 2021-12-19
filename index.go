package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type JsonActivityResponse struct {
	Type    string     `json:"type"`
	Data    []Activity `json:"data"`
	Message string     `json:"message"`
}

type JsonUserResponse struct {
	Type    string `json:"type"`
	Data    User   `json:"data"`
	Message string `json:"message"`
}

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetActivities(w http.ResponseWriter, r *http.Request) {
	printMessage("Getting properties...")

	db := setupDB()

	// Get all properties from table
	rows, err := db.Query("SELECT * FROM activities")

	checkErr(err)

	var activities []Activity

	for rows.Next() {
		var id int
		var t string

		err = rows.Scan(&id, &t)

		checkErr(err)

		activities = append(activities, Activity{ActivityID: id, Type: t})
	}

	var response = JsonActivityResponse{Type: "success", Data: activities}

	json.NewEncoder(w).Encode(response)
}

// create an activity
func CreateActivity(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("type")

	var response = JsonActivityResponse{}

	if t == "" {
		response = JsonActivityResponse{Type: "error", Message: "You are missing a type"}
	} else {
		db := setupDB()

		printMessage("inserting into db")

		var lastInsertID int
		err := db.QueryRow("INSERT INTO activities(type) VALUES($1) returning id;", t).Scan(&lastInsertID)
		// check errors
		checkErr(err)
		fmt.Println(lastInsertID)
		var createdActivity []Activity
		rows, err := db.Query("SELECT * FROM activities WHERE id=($1)", lastInsertID)

		checkErr(err)

		for rows.Next() {
			var id int
			var t string

			err = rows.Scan(&id, &t)

			checkErr(err)

			createdActivity = append(createdActivity, Activity{ActivityID: id, Type: t})
		}
		response = JsonActivityResponse{Type: "success", Data: createdActivity, Message: "The activity has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}

// create a user
func createUser(w http.ResponseWriter, r *http.Request) {
	var response = JsonUserResponse{}
	dt := time.Now().String()
	db := setupDB()

	printMessage("inserting user into db")

	var lastInsertID int
	err := db.QueryRow("INSERT INTO users(createtime) VALUES($1) returning id;", dt).Scan(&lastInsertID)
	// check errors
	checkErr(err)
	fmt.Println(lastInsertID)
	var createdUser User
	rows, err := db.Query("SELECT * FROM users WHERE id=($1)", lastInsertID)

	checkErr(err)

	for rows.Next() {
		var id int
		var time string

		err = rows.Scan(&id, &time)

		checkErr(err)

		createdUser = User{UserID: id, CreatedTime: time}
	}
	response = JsonUserResponse{Type: "success", Data: createdUser, Message: "The user has been inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	// Get all properties
	router.HandleFunc("/activities/", GetActivities).Methods("GET")

	// Create a property
	router.HandleFunc("/activity/", CreateActivity).Methods("POST")

	// Create a user
	router.HandleFunc("/user/", createUser).Methods("POST")

	print("listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
