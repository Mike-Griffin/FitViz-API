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
		var userId int

		err = rows.Scan(&id, &t, &userId)

		checkErr(err)

		activities = append(activities, Activity{ActivityID: id, Type: t, UserID: userId})
	}

	var response = JsonActivityResponse{Type: "success", Data: activities}

	json.NewEncoder(w).Encode(response)
}

// create an activity
func CreateActivity(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("type")
	userId := r.FormValue("user_id")

	var response = JsonActivityResponse{}

	if t == "" {
		response = JsonActivityResponse{Type: "error", Message: "You are missing a type"}
	} else if userId == "" {
		response = JsonActivityResponse{Type: "error", Message: "You are missing a user id"}
	} else {
		db := setupDB()

		printMessage("inserting into db")

		var lastInsertID int
		err := db.QueryRow("INSERT INTO activities(type, user_id) VALUES($1, $2) returning id;", t, userId).Scan(&lastInsertID)
		// check errors
		checkErr(err)
		fmt.Println(lastInsertID)
		var createdActivity []Activity
		rows, err := db.Query("SELECT * FROM activities WHERE id=($1)", lastInsertID)

		checkErr(err)

		for rows.Next() {
			var id int
			var t string
			var userId int

			err = rows.Scan(&id, &t, &userId)

			checkErr(err)

			createdActivity = append(createdActivity, Activity{ActivityID: id, Type: t, UserID: userId})
		}
		response = JsonActivityResponse{Type: "success", Data: createdActivity, Message: "The activity has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}

// delete an activity
func deleteActivity(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	activityID := params["activityid"]

	var response = JsonActivityResponse{}

	if activityID == "" {
		response = JsonActivityResponse{Type: "error", Message: "You are missing activityid parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting activity from DB")

		_, err := db.Exec("DELETE FROM activities where id = $1", activityID)

		// check errors
		checkErr(err)

		response = JsonActivityResponse{Type: "success", Message: "The activity has been deleted successfully!"}
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

// delete a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["userid"]

	var response = JsonUserResponse{}

	if userId == "" {
		response = JsonUserResponse{Type: "error", Message: "You are missing userid parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting user from DB")

		_, err := db.Exec("DELETE FROM users where user_id = $1", userId)

		// check errors
		checkErr(err)

		response = JsonUserResponse{Type: "success", Message: "The user has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()

	// Get all activities
	router.HandleFunc("/activities/", GetActivities).Methods("GET")

	// Create an activity
	router.HandleFunc("/activity/", CreateActivity).Methods("POST")

	// Delete an activity
	router.HandleFunc("/activity/{activityid}", deleteActivity).Methods("DELETE")

	// Create a user
	router.HandleFunc("/user/", createUser).Methods("POST")

	// Delete a user
	router.HandleFunc("/user/{userid}", deleteUser).Methods("DELETE")

	print("listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
