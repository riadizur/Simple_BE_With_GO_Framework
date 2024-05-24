package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"kubota-gasoline-api/internal/models"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFlags(0) // disable default flags
}

func logWithTimestamp(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("%s %s", timestamp, message)
}

func RegisterShiftHandlers(router *mux.Router, db *sql.DB) {
	router.HandleFunc("/api/gasoline/shift-produksi/start", startShift(db)).Methods("POST")
	router.HandleFunc("/api/gasoline/shift-produksi/finish/{id}", finishShift(db)).Methods("PUT")
	router.HandleFunc("/api/gasoline/shift-produksi/list", listShifts(db)).Methods("GET")
	router.HandleFunc("/api/gasoline/shift-produksi/delete/{id}", deleteShift(db)).Methods("DELETE")
}

func startShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Start string `json:"start"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			errorMessage := fmt.Sprintf("Error decoding request: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO shift (date, start, isAlreadyFinish) VALUES (?, ?, ?)")
		if err != nil {
			errorMessage := fmt.Sprintf("Error preparing statement: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := stmt.Exec(request.Start[:10], request.Start, 0)
		if err != nil {
			errorMessage := fmt.Sprintf("Error executing statement: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := res.LastInsertId()
		if err != nil {
			errorMessage := fmt.Sprintf("Error getting last insert ID: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newShift := models.Shift{
			ID:              int(id),
			Date:            request.Start[:10],
			Start:           request.Start,
			Finish:          "",
			IsAlreadyFinish: false,
		}

		broadcast <- newShift

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": true,
			"data":   newShift,
		})
		logWithTimestamp(fmt.Sprintf("New shift started: %+v", newShift))
	}
}

func finishShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorMessage := fmt.Sprintf("Invalid shift ID: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, "Invalid shift ID", http.StatusBadRequest)
			return
		}

		var request struct {
			Finish string `json:"finish"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			errorMessage := fmt.Sprintf("Error decoding request: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("UPDATE shift SET finish = ?, isAlreadyFinish = ? WHERE id = ?")
		if err != nil {
			errorMessage := fmt.Sprintf("Error preparing statement: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := stmt.Exec(request.Finish, 1, id)
		if err != nil {
			errorMessage := fmt.Sprintf("Error executing statement: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			errorMessage := fmt.Sprintf("Error getting rows affected: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Shift not found", http.StatusNotFound)
			errorMessage := fmt.Sprintf("Shift not found: ID %d", id)
			logWithTimestamp(errorMessage)
			return
		}

		updatedShift := models.Shift{
			ID:              id,
			Finish:          request.Finish,
			IsAlreadyFinish: true,
		}

		broadcast <- updatedShift

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": true,
			"data": map[string]int{
				"affectedRows":  int(rowsAffected),
				"insertId":      0,
				"warningStatus": 0,
			},
		})
		logWithTimestamp(fmt.Sprintf("Shift finished: %+v", updatedShift))
	}
}

func listShifts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, date, start, finish, isAlreadyFinish FROM shift")
		if err != nil {
			errorMessage := fmt.Sprintf("Error querying shifts: %s", err.Error())
			logWithTimestamp(errorMessage)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var shifts []models.Shift
		for rows.Next() {
			var shift models.Shift
			var isAlreadyFinish int
			if err := rows.Scan(&shift.ID, &shift.Date, &shift.Start, &shift.Finish, &isAlreadyFinish); err != nil {
				errorMessage := fmt.Sprintf("Error scanning row: %s", err.Error())
				logWithTimestamp(errorMessage)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			shift.IsAlreadyFinish = isAlreadyFinish == 1
			shifts = append(shifts, shift)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": true,
			"data":   shifts,
		})
		logWithTimestamp(fmt.Sprintf("Shifts listed: %+v", shifts))
	}
}

func deleteShift(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			logWithTimestamp("Invalid shift ID: " + err.Error())
			http.Error(w, "Invalid shift ID", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("DELETE FROM shift WHERE id = ?")
		if err != nil {
			logWithTimestamp("Error preparing statement: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := stmt.Exec(id)
		if err != nil {
			logWithTimestamp("Error executing statement: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			logWithTimestamp("Error getting rows affected: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Shift not found", http.StatusNotFound)
			logWithTimestamp("Shift not found: ID " + strconv.Itoa(id))
			return
		}

		broadcast <- map[string]int{"deletedId": id}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": true,
			"data": map[string]int{
				"affectedRows": int(rowsAffected),
			},
		})
		logWithTimestamp("Shift deleted: ID " + strconv.Itoa(id))
	}
}
