package main

import (
	"Task_Manager/config"
	"Task_Manager/handler/task"
	"Task_Manager/handler/user"
	Task2 "Task_Manager/service/task"
	User2 "Task_Manager/service/user"
	Task3 "Task_Manager/store/task"
	User3 "Task_Manager/store/user"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	config.DataBaseConfig()
	db := config.DB
	// Init user dependencies
	userStore := User3.NewUserStore(db)
	userService := User2.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)
	// Init task dependencies
	taskStore := Task3.NewStore(db)
	taskService := Task2.NewService(taskStore, userService)
	taskHandler := task.NewHandler(taskService)
	// Setup router
	r := mux.NewRouter()
	// Task routes
	r.HandleFunc("/task", taskHandler.Create).Methods("POST")
	r.HandleFunc("/task/{id}", taskHandler.GetTask).Methods("GET")
	r.HandleFunc("/task/{id}", taskHandler.Complete).Methods("PUT")
	r.HandleFunc("/task/{id}", taskHandler.Delete).Methods("DELETE")
	r.HandleFunc("/task", taskHandler.All).Methods("GET")
	r.HandleFunc("/task/user/{userid}", taskHandler.GetTasksByUserID).Methods("GET")
	// User routes
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	fmt.Println("Server running at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

//Create task table
//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
//	id INT AUTO_INCREMENT PRIMARY KEY,
//	description VARCHAR(255),
//	status BOOLEAN ,
//  userid INT UNSIGNED NOT NULL
//)`)
//if err != nil {
//	log.Fatal("Failed to create tasks table:", err)
//}
//
//// Create user table
//_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
//	id INT AUTO_INCREMENT PRIMARY KEY,
//	name VARCHAR(100),
//	email VARCHAR(100)
//)`)
//if err != nil {
//	log.Fatal("Failed to create users table:", err)
//}
