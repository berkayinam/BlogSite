package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamHandler struct {
	repo *TeamRepository
}

func NewTeamHandler(repo *TeamRepository) *TeamHandler {
	return &TeamHandler{repo: repo}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	team.Members = []TeamMember{{
		Username: username,
		Role:     "admin",
	}}

	if err := h.repo.CreateTeam(r.Context(), &team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := h.repo.GetTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	teams, err := h.repo.GetUserTeams(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	team, err := h.repo.GetTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if user is admin
	isAdmin := false
	for _, member := range team.Members {
		if member.Username == username && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var updatedTeam Team
	if err := json.NewDecoder(r.Body).Decode(&updatedTeam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTeam.ID = id
	updatedTeam.Members = team.Members // Preserve existing members

	if err := h.repo.UpdateTeam(r.Context(), &updatedTeam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeam)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	team, err := h.repo.GetTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if user is admin
	isAdmin := false
	for _, member := range team.Members {
		if member.Username == username && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.repo.DeleteTeam(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	var invite TeamInvite
	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	team, err := h.repo.GetTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if user is admin
	isAdmin := false
	for _, member := range team.Members {
		if member.Username == username && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Check if user is already a member
	for _, member := range team.Members {
		if member.Username == invite.Username {
			http.Error(w, "User is already a member", http.StatusBadRequest)
			return
		}
	}

	invite.TeamID = id
	if err := h.repo.CreateInvite(r.Context(), &invite); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invite)
}

func (h *TeamHandler) RespondToInvite(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid invite ID", http.StatusBadRequest)
		return
	}

	var response struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if response.Status != "accepted" && response.Status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	invite, err := h.repo.GetInvite(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if invite.Username != username {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.repo.UpdateInviteStatus(r.Context(), id, response.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response.Status == "accepted" {
		member := TeamMember{
			Username: username,
			Role:     "member",
		}
		if err := h.repo.AddMember(r.Context(), invite.TeamID, member); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	memberUsername := vars["username"]
	username := r.Header.Get("username")

	team, err := h.repo.GetTeam(r.Context(), teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if user is admin
	isAdmin := false
	for _, member := range team.Members {
		if member.Username == username && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Check if target member is admin
	for _, member := range team.Members {
		if member.Username == memberUsername && member.Role == "admin" {
			http.Error(w, "Cannot remove admin member", http.StatusForbidden)
			return
		}
	}

	if err := h.repo.RemoveMember(r.Context(), teamID, memberUsername); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 