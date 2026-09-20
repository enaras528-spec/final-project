package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"todo_app/pkg/db"
)

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ошибка кодирования JSON: %v", err)
	}
}
func sendError(w http.ResponseWriter, errText string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	errResponse := map[string]string{"error": errText}
	json.NewEncoder(w).Encode(errResponse)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendError(w, "Не указан заголовок", http.StatusBadRequest)
		return
	}
	now := time.Now()
	// Оставляем от текущего времени только дату
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	todayStr := today.Format(DateFormat)

	if task.Date == "" {
		task.Date = todayStr
	}

	parsedDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		sendError(w, "Дата представлена в формате, отличном от 20060102", http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		next, err := NextDate(today, task.Date, task.Repeat)
		if err != nil {
			sendError(w, "Правило повторения указано в неправильном формате", http.StatusBadRequest)
			return
		}
		if parsedDate.Before(today) {
			task.Date = next
		}
	} else {
		if parsedDate.Before(today) {
			task.Date = todayStr
		}
	}

	id, err := db.AddTask(task)
	if err != nil {
		sendError(w, "Ошибка при добавлении задачи в БД", http.StatusInternalServerError)
		return
	}
	sendJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}
