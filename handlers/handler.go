package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service APIServingService
}

func (h *Handler) initHandler(service APIServingService) *Handler {
	h.service = service
	return h
}

func NewHandler(service APIServingService) *Handler {
	return (&Handler{}).initHandler(service)
}

func (h *Handler) HandleLastHour(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.service.GetLastHourExpenses()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(expenses)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleRange(w http.ResponseWriter, r *http.Request) {
	startHourParam := r.URL.Query().Get("startHour")
	if len(startHourParam) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	startHour, err := strconv.Atoi(startHourParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	normalizedStartHour := time.Unix(int64(startHour), 0).Truncate(time.Hour).Unix()
	endHourParam := r.URL.Query().Get("endHour")
	if len(endHourParam) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	endHour, err := strconv.Atoi(endHourParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	normalizedEndHour := time.Unix(int64(endHour), 0).Truncate(time.Hour).Unix()

	if normalizedEndHour < normalizedStartHour {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expenses, err := h.service.GetExpenses(normalizedStartHour, normalizedEndHour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(expenses)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
