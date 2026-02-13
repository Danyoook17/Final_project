package api

import (
	"net/http"
	"todoapp/pkg/db"
)

const DefaultTasksLimit = 50

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(DefaultTasksLimit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, tasksResp{Tasks: tasks})
}
