package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ToDoListItem struct {
	Id          int
	Description string
	Date        string
	Creator     string
	Completed   bool
	Editing     bool
}

var db *sql.DB
var err error

func main() {

	fmt.Println("Connecting to database...")
	connString := "sql3335274:R2PtXLpGNU@tcp(sql3.freemysqlhosting.net:3306)/sql3335274"
	db, err = sql.Open("mysql", connString)

	if err != nil {
		panic(err.Error())
	} else {
		log.Println("Connection Established")
	}
	defer db.Close()

	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	router.HandleFunc("/api/getToDoListItems", getToDoListItems).Methods("GET")
	router.HandleFunc("/api/createToDoListItem", createToDoListItem).Methods("POST")
	router.HandleFunc("/api/getToDoListItemById/{id}", getToDoListItemByID).Methods("GET")
	router.HandleFunc("/api/updateTodoListItem/{id}", updateTodoListItem).Methods("PUT")
	router.HandleFunc("/api/deleteToDoListItem/{id}", deleteToDoListItem).Methods("DELETE")

	http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(router))
}

func getToDoListItems(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var toDoListItems []ToDoListItem
	// Enable cors
	enableCors(&w)

	result, err := db.Query("SELECT * FROM 	`ToDoListItemInfo`")
	//fmt.Println("Reacheddddd")

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()
	// Looping through the list Items
	for result.Next() {
		var toDoItem ToDoListItem
		err := result.Scan(&toDoItem.Id, &toDoItem.Description, &toDoItem.Date, &toDoItem.Creator, &toDoItem.Completed, &toDoItem.Editing)
		if err != nil {
			panic(err.Error())
		}
		toDoListItems = append(toDoListItems, toDoItem)
	}
	json.NewEncoder(w).Encode(toDoListItems)

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func createToDoListItem(w http.ResponseWriter, r *http.Request) {
	// Enable cors
	enableCors(&w)
	stmt, err := db.Prepare("INSERT INTO `ToDoListItemInfo`(Description, Date, Creator, Completed, editing) VALUES( ?,?,?,?,?)")

	if err != nil {
		panic(err.Error())
	}
	// Readgin the data from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("I hittt crerateeeee")
	// Extracting the different todolist item values.
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	description := keyVal["Description"]
	date := keyVal["Date"]
	creator := keyVal["Creator"]
	completed := keyVal["Completed"]
	editing := keyVal["editing"]
	fmt.Println("reached --shabih")
	_, err = stmt.Exec(description, date, creator, completed, editing)

	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "New Todo Item was created")
}

// Although this api isn't used , i though to add it for futurte use.
func getToDoListItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT  * FROM `ToDoListItemInfo` WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var toDoListItem ToDoListItem
	for result.Next() {
		err := result.Scan(&toDoListItem.Id, &toDoListItem.Description, &toDoListItem.Date, &toDoListItem.Creator)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(toDoListItem)
}

func updateTodoListItem(w http.ResponseWriter, r *http.Request) {

	// Enable cors
	enableCors(&w)
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE `ToDoListItemInfo` SET Description = ? , Date = ? , Creator = ? , Completed = ? , Editing = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newdescription := keyVal["Description"]
	newDate := keyVal["Date"]
	newCreator := keyVal["Creator"]
	newCompleted, err := strconv.ParseBool(keyVal["Completed"])
	newEditing, err := strconv.ParseBool(keyVal["Editing"])
	_, err = stmt.Exec(newdescription, newDate, newCreator, newCompleted, newEditing, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with Id = %s was updated", params["id"])
}

func deleteToDoListItem(w http.ResponseWriter, r *http.Request) {
	// Enable cors
	enableCors(&w)
	fmt.Println("Reached inside delete module")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM `ToDoListItemInfo` WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with Id = %s was deleted", params["id"])
}
