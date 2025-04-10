package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "auth healty")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Şifreyi hashle
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	// MongoDB bağlantısı
	collection := Client.Database("authdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Aynı kullanıcı var mı kontrol et
	count, _ := collection.CountDocuments(ctx, map[string]interface{}{"username": user.Username})
	if count > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqUser User
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil || reqUser.Username == "" || reqUser.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Mongo bağlantısı
	collection := Client.Database("authdb").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storedUser User
	err = collection.FindOne(ctx, map[string]interface{}{"username": reqUser.Username}).Decode(&storedUser)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(reqUser.Password, storedUser.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// JWT oluştur
	token, err := GenerateJWT(storedUser.Username)
	if err != nil {
		http.Error(w, "JWT error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
