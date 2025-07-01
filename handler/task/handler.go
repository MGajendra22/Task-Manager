package task

import (
	"Task_Manager/model/task"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	svc TaskServiceInterface
}

// NewHandler : Factory function to implement and return behaviour
func NewHandler(s TaskServiceInterface) *Handler {
	return &Handler{svc: s}
}

// Create Task (POST /task)
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	var t task.Task

	if err = json.Unmarshal(body, &t); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	task1, err := h.svc.Create(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task1)

	if err != nil {
		http.Error(w, "Failed to marshal task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(resp); err != nil {
		fmt.Println("Write failed:", err)
	}

}

// GetTask by ID (GET /task/{id})
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := mux.Vars(r)["id"] // same as r.pathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task1, err := h.svc.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(task1)

	if err != nil {
		http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		fmt.Println("Write failed:", err)
	}
}

// GetTasksByUserID which are assigned to user_id
func (h *Handler) GetTasksByUserID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := mux.Vars(r)["userid"]

	userid, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	tasks, err := h.svc.GetTasksByUserID(userid)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Complete Task (PUT /task/{id})
func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err = h.svc.Complete(id); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(fmt.Sprintf("Task %d marked as complete", id)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete Task (DELETE /task/{id})
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err = h.svc.Delete(id); err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(fmt.Sprintf("Task %d deleted", id)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// All Tasks (GET /task)
func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.svc.All()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, "Failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(resp); err != nil {
		http.Error(w, "Write failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
