package handlers

import (
	dao "GOLANG/project/dal"
	"GOLANG/project/models"
	"GOLANG/project/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var userDaoInstance *dao.UserDao

func UserRoutes(router *mux.Router) {
	userDaoInstance, _ = dao.InitializeUserDB()
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/signup", signUpHandler).Methods("POST")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	ctx := context.Background()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	savedPwd, err := userDaoInstance.GetPassword(ctx, user)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*savedPwd), []byte(user.Pwd))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		http.Error(w, "Could not create authenticate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	ctx := context.Background()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user.HashPassword()
	err = userDaoInstance.AddUser(ctx, user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully", user.Email)
}
