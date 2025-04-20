package internal

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentHandler struct {
	commentRepo *CommentRepository
	postLikeRepo *PostLikeRepository
	friendshipRepo *FriendshipRepository
}

func NewCommentHandler(cr *CommentRepository, pr *PostLikeRepository, fr *FriendshipRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo: cr,
		postLikeRepo: pr,
		friendshipRepo: fr,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	comment.PostID = vars["postId"]
	comment.Author = r.Header.Get("username")

	if err := h.commentRepo.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	comments, err := h.commentRepo.GetCommentsByPost(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var update struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.commentRepo.UpdateComment(r.Context(), commentID, update.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := h.commentRepo.DeleteComment(r.Context(), commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	if err := h.commentRepo.LikeComment(r.Context(), commentID, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) UnlikeComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	username := r.Header.Get("username")
	if err := h.commentRepo.UnlikeComment(r.Context(), commentID, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Post Like handlers
func (h *CommentHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	username := r.Header.Get("username")

	like := &PostLike{
		PostID: postID,
		Username: username,
	}

	if err := h.postLikeRepo.LikePost(r.Context(), like); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]
	username := r.Header.Get("username")

	if err := h.postLikeRepo.UnlikePost(r.Context(), postID, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) GetPostLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	likes, err := h.postLikeRepo.GetPostLikes(r.Context(), postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(likes)
}

// Friendship handlers
func (h *CommentHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var friendship Friendship
	if err := json.NewDecoder(r.Body).Decode(&friendship); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	friendship.User1 = r.Header.Get("username")

	if err := h.friendshipRepo.SendFriendRequest(r.Context(), &friendship); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	requests, err := h.friendshipRepo.GetFriendRequests(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

func (h *CommentHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	if err := h.friendshipRepo.AcceptFriendRequest(r.Context(), requestID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CommentHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")

	friends, err := h.friendshipRepo.GetFriends(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
} 