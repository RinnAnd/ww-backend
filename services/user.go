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
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&login)

	rows, err := us.sql.Query("SELECT * FROM users WHERE email = $1", login.Email)
	if err != nil {
		json.NewEncoder(w).Encode("There was an error fetching the user")
		return
	}

	if rows.Next() {
		err := rows.Scan(&user.ID, &user.UserName, &user.Name, &user.Email, &user.Password)
		if err != nil {
			json.NewEncoder(w).Encode("There was an error fetching the password")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	success := utils.CheckPasswordHash(login.Password, user.Password)

	if !success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Password is incorrect",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusAccepted,
		Message: "User session accepted",
		Data: models.Session{
			ID:       user.ID,
			Name:     user.Name,
			UserName: user.UserName,
			Email:    user.Email,
		},
	})
}

func (us *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}

	rows, err := us.sql.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, "There was an error fetching users", http.StatusBadRequest)
		return
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.Name, &user.Email, &user.Password)
		if err != nil {
			http.Error(w, "There was an error fetching user", http.StatusBadRequest)
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusOK,
		Message: "Users fetched",
		Data:    users,
	})
}

func (us *UserService) ChangePassword(w http.ResponseWriter, r *http.Request) {
	emailPassword := make(map[string]string)

	json.NewDecoder(r.Body).Decode(&emailPassword)

	newPassword := emailPassword["password"]

	err := utils.HashPassword(&newPassword)
	if err != nil {
		http.Error(w, "There was an error hashing password", http.StatusInternalServerError)
	}

	_, err = us.sql.Exec("UPDATE users SET password = $1 WHERE email = $2", newPassword, emailPassword["email"])
	if err != nil {
		http.Error(w, "There was an error updating the password", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusOK,
		Message: "The user password has been updated successfully",
		Data:    nil,
	})
}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {}

func MakeUserService(db *sql.DB) *UserService {
	return &UserService{
		sql: db,
	}
}
