package api

import (
	"net/http"

	"todo_app/pkg/db"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	const limit = 50
	search := r.FormValue("search")

	tasks, err := db.Tasks(search, limit)
	if err != nil {
		sendError(w, "Ошибка при получении списка задач", http.StatusBadRequest)
		return
	}
	sendJSON(w, http.StatusOK, map[string]any{
		"tasks": tasks,
	})
}
