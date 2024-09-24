package controllers

import (
	"Test/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB // Biến toàn cục cho database

func InitController(database *sql.DB) {
	db = database
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy các tên cột từ query params (ví dụ: ?columns=id,name)
	// columnsParam := r.URL.Query().Get("columns")
	columnsParam := "name,id"
	log.Println(columnsParam)
	if columnsParam == "" {
		http.Error(w, "Missing columns parameter", http.StatusBadRequest)
		return
	}
	columns := strings.Split(columnsParam, ",")
	method := models.NewMethod(db, "users")
	results, err := method.GetAllRecord(columns, "WHERE id = 8")

	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, results)
}

// Controller để lấy thông tin bản ghi (GET)
func FindUserHandler(w http.ResponseWriter, r *http.Request) {
	recordIDStr := r.URL.Query().Get("id")
	if recordIDStr == "" {
		http.Error(w, "Missing record ID", http.StatusBadRequest)
		return
	}
	recordID, _ := strconv.Atoi(recordIDStr)
	method := models.NewMethod(db, "users")
	record, err := method.GetRecordByID(recordID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, record)
}

// Controller để thêm bản ghi (POST)
func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fields := []string{"name"}
	values := []interface{}{data["name"]}
	method := models.NewMethod(db, "users")
	if err := method.InsertRecord(fields, values); err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Record added successfully"})
}

// Controller để cập nhật bản ghi (PUT)
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	recordIDStr := r.URL.Query().Get("id")
	if recordIDStr == "" {
		http.Error(w, "Missing record ID", http.StatusBadRequest)
		return
	}
	recordID, _ := strconv.Atoi(recordIDStr)
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fields := []string{"name"}
	values := []interface{}{data["name"]}
	method := models.NewMethod(db, "users")
	if err := method.UpdateRecord(recordID, fields, values); err != nil {
		http.Error(w, "Update error", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Record updated successfully"})
}

// Controller để xóa bản ghi (DELETE)
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	recordIDStr := r.URL.Query().Get("id")
	if recordIDStr == "" {
		http.Error(w, "Missing record ID", http.StatusBadRequest)
		return
	}
	recordID, _ := strconv.Atoi(recordIDStr)
	method := models.NewMethod(db, "users")
	if err := method.DeleteRecord(recordID); err != nil {
		http.Error(w, "Delete error", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Record deleted successfully"})
}
