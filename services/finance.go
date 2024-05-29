package services

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/RinnAnd/ww-backend/database"
	"github.com/RinnAnd/ww-backend/models"
	"github.com/RinnAnd/ww-backend/utils"
)

type FinanceService struct {
	db database.Database
}

func NewFinanceService(db database.Database) *FinanceService {
	return &FinanceService{
		db: db,
	}
}

func (fs *FinanceService) CreateFinance(w http.ResponseWriter, r *http.Request) {
	finance := models.Finance{}
	json.NewDecoder(r.Body).Decode(&finance)

	err := fs.db.CreateFinance(&finance)
	if err != nil {
		http.Error(w, "There was an error creating the finance", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response[models.Finance]{
		Status:  http.StatusAccepted,
		Message: "Finance added successfully",
		Data:    finance,
	})
}

func (fs *FinanceService) GetUserSavings(id string, finance *models.RelFinance) error {
	savings, err := fs.db.GetSavingsByFinanceID(id)
	if err != nil {
		return err
	}

	finance.Savings = savings

	return nil
}

func (fs *FinanceService) GetUserExpenses(id string, finance *models.RelFinance) error {
	expenses, err := fs.db.GetExpensesByFinanceID(id)
	if err != nil {
		return err
	}

	finance.Expenses = expenses

	return nil
}

func (fs *FinanceService) GetUserFinances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	finances, err := fs.db.GetFinancesByUserID(userId)
	if err != nil {
		http.Error(w, "There was an error fetching the finances", http.StatusInternalServerError)
		return
	}

	if len(finances) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.Response[any]{
			Status:  http.StatusNotFound,
			Message: "No finances found",
			Data:    nil,
		})
		return
	}

	for i := range finances {
		var wg sync.WaitGroup
		wg.Add(2)

		errChan := make(chan error, 2)

		go func(i int) {
			defer wg.Done()
			err := fs.GetUserExpenses(finances[i].ID, finances[i])
			if err != nil {
				errChan <- err
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			err := fs.GetUserSavings(finances[i].ID, finances[i])
			if err != nil {
				errChan <- err
			}
		}(i)

		wg.Wait()

		close(errChan)

		if len(errChan) > 0 {
			http.Error(w, "There was an error fetching savings or expenses", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response[[]*models.RelFinance]{
		Status:  http.StatusOK,
		Message: "Retrieved user finances",
		Data:    finances,
	})
}

func (fs *FinanceService) CreateExpense(w http.ResponseWriter, r *http.Request) {
	expense := models.Expense{}

	json.NewDecoder(r.Body).Decode(&expense)

	err := fs.db.CreateExpense(&expense)
	if err != nil {
		http.Error(w, "There was an error creating the expense", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response[models.Expense]{
		Status:  http.StatusAccepted,
		Message: "Created the expense successfully",
		Data:    expense,
	})
}

func (fs *FinanceService) CreateSaving(w http.ResponseWriter, r *http.Request) {
	saving := models.Saving{}

	json.NewDecoder(r.Body).Decode(&saving)

	err := fs.db.CreateSaving(&saving)
	if err != nil {
		http.Error(w, "There was an error creating the saving", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response[models.Saving]{
		Status:  http.StatusAccepted,
		Message: "Created the saving successfully",
		Data:    saving,
	})
}
