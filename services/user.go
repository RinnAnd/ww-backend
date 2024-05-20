package services

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/RinnAnd/ww-backend/models"
	"github.com/RinnAnd/ww-backend/utils"
)

type UserService struct {
	sql *sql.DB
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&user)

	err := utils.HashPassword(&user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = us.sql.Exec("INSERT INTO users (username, name, email, password) VALUES ($1, $2, $3, $4)", user.UserName, user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "There was an error inserting user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusCreated,
		Message: "User created",
		Data:    user,
	})
}

func (us *UserService) Auth(w http.ResponseWriter, r *http.Request) {
	login := models.Login{}
	var password string

	json.NewDecoder(r.Body).Decode(&login)

	rows, err := us.sql.Query("SELECT password FROM users WHERE email = $1", login.Email)
	if err != nil {
		json.NewEncoder(w).Encode("There was an error fetching the user")
		return
	}

	if rows.Next() {
		err := rows.Scan(&password)
		if err != nil {
			json.NewEncoder(w).Encode("There was an error fetching the password")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("User not found")
		return
	}

	success := utils.CheckPasswordHash(login.Password, password)

	//! Return a user auth with a token
	json.NewEncoder(w).Encode(success)
}

func (us *UserService) GetUser(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) EditUser(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {}

func MakeUserService(db *sql.DB) *UserService {
	return &UserService{
		sql: db,
	}
}
