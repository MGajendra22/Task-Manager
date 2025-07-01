package user

import (
	"Task_Manager/model/user"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type UserHandler struct {
	Service UserServiceInterface
}

// NewUserHandler : Factory function to implement and return behaviour
func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{Service: service}
}

// CreateUser : Creates user (POST /task)
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read entire request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close body")
		}
	}(r.Body)
	// Unmarshal into struct
	var user1 user.User

	if err = json.Unmarshal(body, &user1); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}
	// Validate and save
	createdUser, err := h.Service.Create(user1)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with created user
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// GetUser : To retrieve user with user-id
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user1, err := h.Service.Get(id)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusNotFound)
		return
	}

	resp, _ := json.Marshal(user1) // Convert struct to JSON
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteUser : To delete user with user-id
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err = h.Service.Delete(id); err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(fmt.Sprintf("User %d Removed", id)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetAllUsers : To retrieve all users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	users, err := h.Service.All()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(users) // Convert struct to JSON
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
