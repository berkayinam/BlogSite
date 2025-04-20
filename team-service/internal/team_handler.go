package internal

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	repo *TeamRepository
}

func NewTeamHandler(repo *TeamRepository) *TeamHandler {
	return &TeamHandler{repo: repo}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get username from context
	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	team, err := h.repo.CreateTeam(request.Name, request.Description, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := h.repo.GetTeam(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetUserTeams(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teams, err := h.repo.GetUserTeams(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin, err := h.repo.IsMemberAdmin(teamID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only team admins can update team details", http.StatusForbidden)
		return
	}

	if err := h.repo.UpdateTeam(teamID, request.Name, request.Description); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin, err := h.repo.IsMemberAdmin(teamID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only team admins can delete teams", http.StatusForbidden)
		return
	}

	if err := h.repo.DeleteTeam(teamID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TeamID          int    `json:"team_id"`
		InviteeUsername string `json:"invitee_username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	inviterUsername := r.Context().Value("username").(string)
	if inviterUsername == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin, err := h.repo.IsMemberAdmin(request.TeamID, inviterUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only team admins can invite members", http.StatusForbidden)
		return
	}

	// Check if user is already a member
	isMember, err := h.repo.IsMember(request.TeamID, request.InviteeUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isMember {
		http.Error(w, "User is already a team member", http.StatusBadRequest)
		return
	}

	invite, err := h.repo.CreateInvite(request.TeamID, inviterUsername, request.InviteeUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invite)
}

func (h *TeamHandler) RespondToInvite(w http.ResponseWriter, r *http.Request) {
	var request struct {
		InviteID int    `json:"invite_id"`
		Status   string `json:"status"` // "accepted" or "rejected"
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invite, err := h.repo.GetInvite(request.InviteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if invite.InviteeUsername != username {
		http.Error(w, "This invite is not for you", http.StatusForbidden)
		return
	}

	if err := h.repo.UpdateInviteStatus(request.InviteID, request.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Status == "accepted" {
		if err := h.repo.AddMember(invite.TeamID, username); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.Atoi(vars["teamId"])
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	memberUsername := vars["username"]
	if memberUsername == "" {
		http.Error(w, "Member username is required", http.StatusBadRequest)
		return
	}

	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	isAdmin, err := h.repo.IsMemberAdmin(teamID, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAdmin {
		http.Error(w, "Only team admins can remove members", http.StatusForbidden)
		return
	}

	if err := h.repo.RemoveMember(teamID, memberUsername); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetUserInvites(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	if username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	invites, err := h.repo.GetUserInvites(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(invites)
} 