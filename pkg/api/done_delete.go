package api

import (
	"net/http"
	"strings"
	"time"

	"todoapp/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})
		return
	}

	now := dateOnly(time.Now())
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}
