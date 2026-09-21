package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	store *TaskStore
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusUnprocessableEntity, "title is required")
		return
	}

	task, err := h.store.Create(r.Context(), req.Title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func taskID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	id, err := taskID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return
	}

	task, err := h.store.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := taskID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return
	}

	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Title != nil {
		v := strings.TrimSpace(*req.Title)
		if v == "" {
			writeError(w, http.StatusUnprocessableEntity, "title is required")
			return
		}
		req.Title = &v
	}
	if req.Title == nil && req.Done == nil {
		writeError(w, http.StatusUnprocessableEntity, "at least one field is required")
		return
	}

	task, err := h.store.Update(r.Context(), id, req.Title, req.Done)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	id, err := taskID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id must be a number")
		return
	}
	if err := h.store.Delete(r.Context(), id); errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
