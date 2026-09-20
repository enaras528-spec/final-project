package api

import (
	"encoding/json"
	"net/http"
	"time"

	"todo_app/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	sendJSON(w, http.StatusOK, task)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		sendError(w, "Не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		sendError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now()
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

	err = db.UpdateTask(task)
	if err != nil {
		sendError(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	sendJSON(w, http.StatusOK, map[string]any{})
}
