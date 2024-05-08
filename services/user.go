package services

import (
	"encoding/json"
	"net/http"

	"github.com/RinnAnd/ww-backend/models"
)

type UserService struct {
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	json.NewEncoder(w).Encode(user)
}

func (us *UserService) GetUser(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) EditUser(w http.ResponseWriter, r *http.Request) {}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {}

func MakeUserService() *UserService {
	return &UserService{}
}
