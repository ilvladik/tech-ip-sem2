package order

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	repo   *Repo
	client *UserServiceClient
}

func NewHandler(repo *Repo, client *UserServiceClient) *Handler {
	return &Handler{
		repo:   repo,
		client: client,
	}
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		writeOrderNotFound(w, id)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) GetOrderWithUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.repo.GetByID(id)
	if err != nil {
		writeOrderNotFound(w, id)
		return
	}

	user, err := h.client.GetUserByID(order.UserID)
	if err != nil {
		http.Error(w, "failed to get user data: "+err.Error(), http.StatusBadGateway)
		return
	}

	result := OrderWithUser{
		Order: order,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	orders := h.repo.GetByUserID(userID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(orders)
}

func writeOrderNotFound(w http.ResponseWriter, id int64) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": "заказ с таким id не найден в системе",
		"id":    id,
	})
}
