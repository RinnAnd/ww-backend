package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ww-backend/models"
)

type FinanceService struct {
}

func (fs *FinanceService) CreateService(w http.ResponseWriter, r *http.Request) {
	finance := models.Finance{}
	json.NewDecoder(r.Body).Decode(&finance)
	fmt.Println(finance)
}
