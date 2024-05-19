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

	_, err = us.sql.Exec("INSERT INTO users VALUES (username = $1, name = $2, email = $3, password = $4)", user.UserName, user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "There was an error inserting user", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(utils.Response{
		Status:  200,
		Message: "User created",
		Data:    user,
	})
}

func (us *UserService) Auth(w http.ResponseWriter, r *http.Request) {
	user := models.Login{}

	json.NewDecoder(r.Body).Decode(&user)

	success := utils.CheckPasswordHash(user.Password, "")

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
