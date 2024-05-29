package services

import (
	"encoding/json"
	"net/http"

	"github.com/RinnAnd/ww-backend/database"
	"github.com/RinnAnd/ww-backend/models"
	"github.com/RinnAnd/ww-backend/utils"
)

type UserService struct {
	db database.Database
}

func NewUserService(db database.Database) *UserService {
	return &UserService{
		db: db,
	}
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&user)

	err := utils.HashPassword(&user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err = us.db.CreateUser(&user)
	if err != nil {
		http.Error(w, "There was an error creating the user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.Response[models.User]{
		Status:  http.StatusCreated,
		Message: "User created",
		Data:    user,
	})
}

func (us *UserService) Auth(w http.ResponseWriter, r *http.Request) {
	login := models.Login{}

	json.NewDecoder(r.Body).Decode(&login)

	user, err := us.db.GetUserByEmail(login.Email)
	if err != nil {
		http.Error(w, "There was an error fetching the user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.Response[any]{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	success := utils.CheckPasswordHash(login.Password, user.Password)

	if !success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response[any]{
			Status:  http.StatusBadRequest,
			Message: "Password is incorrect",
			Data:    nil,
		})
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response[models.Session]{
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
	users, err := us.db.GetAllUsers()
	if err != nil {
		http.Error(w, "There was an error fetching the users", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response[[]*models.User]{
		Status:  http.StatusOK,
		Message: "Users fetched",
		Data:    users,
	})
}

func (us *UserService) ChangePassword(w http.ResponseWriter, r *http.Request) {
	emailPassword := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	json.NewDecoder(r.Body).Decode(&emailPassword)

	newPassword := emailPassword.Password

	err := utils.HashPassword(&newPassword)
	if err != nil {
		http.Error(w, "There was an error hashing password", http.StatusInternalServerError)
	}

	err = us.db.UpdatePasswordForEmail(emailPassword.Email, newPassword)
	if err != nil {
		http.Error(w, "There was an error updating the password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response[any]{
		Status:  http.StatusOK,
		Message: "The user password has been updated successfully",
		Data:    nil,
	})
}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {}
