package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/RinnAnd/ww-backend/models"
	"github.com/RinnAnd/ww-backend/utils"
	"github.com/gorilla/mux"
)

type FinanceService struct {
	sql *sql.DB
}

type Finance struct {
	ID       string
	UserID   string
	Month    int
	Year     int
	Salary   int
	Expenses []models.Expense
	Savings  []models.Saving
}

func (fs *FinanceService) CreateFinance(w http.ResponseWriter, r *http.Request) {
	finance := models.Finance{}
	json.NewDecoder(r.Body).Decode(&finance)

	_, err := fs.sql.Exec("INSERT INTO finances (user_id, month, year, salary) VALUES ($1, $2, $3, $4)", finance.UserID, finance.Month, finance.Year, finance.Salary)
	if err != nil {
		http.Error(w, "There was an error creating finance", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusAccepted,
		Message: "Finance added successfully",
		Data:    finance,
	})
}

func (fs *FinanceService) GetUserSavings(id string, finance *Finance) error {
	savings := []models.Saving{}

	rows, err := fs.sql.Query("SELECT * FROM savings WHERE finance_id = $1", id)
	if err != nil {
		return err
	}

	for rows.Next() {
		saving := models.Saving{}
		err := rows.Scan(&saving.ID, &saving.Finance, &saving.Amount)
		if err != nil {
			return err
		}

		savings = append(savings, saving)
	}

	finance.Savings = savings

	return nil
}

func (fs *FinanceService) GetUserExpenses(id string, finance *Finance) error {
	expenses := []models.Expense{}

	rows, err := fs.sql.Query("SELECT * FROM expenses WHERE finance_id = $1", id)
	if err != nil {
		fmt.Println("Error in the query")
		return err
	}

	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Finance, &expense.Name, &expense.Amount, &expense.Category)
		if err != nil {
			return err
		}

		expenses = append(expenses, expense)
	}

	finance.Expenses = expenses

	return nil
}

func (fs *FinanceService) GetUserFinances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	finances := []Finance{}

	rows, err := fs.sql.Query(`SELECT * FROM finances WHERE user_id = $1`, userId)
	if err != nil {
		http.Error(w, "There was an error fetching the user's finances", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		finance := Finance{}
		err := rows.Scan(&finance.ID, &finance.UserID, &finance.Month, &finance.Year, &finance.Salary)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "There was an error fetching finance", http.StatusBadRequest)
			return
		}

		finances = append(finances, finance)
	}

	if len(finances) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.Response{
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
			err := fs.GetUserExpenses(finances[i].ID, &finances[i])
			if err != nil {
				errChan <- err
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			err := fs.GetUserSavings(finances[i].ID, &finances[i])
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
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusOK,
		Message: "Retrieved user finances",
		Data:    finances,
	})
}

func (fs *FinanceService) CreateExpense(w http.ResponseWriter, r *http.Request) {
	expense := models.Expense{}

	json.NewDecoder(r.Body).Decode(&expense)

	_, err := fs.sql.Exec("INSERT INTO expenses (finance_id, name, amount, category) VALUES ($1, $2, $3, $4)", expense.Finance, expense.Name, expense.Amount, expense.Category)
	if err != nil {
		http.Error(w, "There was an error inserting the expense", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusAccepted,
		Message: "Created the expense successfully",
		Data:    expense,
	})
}

func (fs *FinanceService) CreateSaving(w http.ResponseWriter, r *http.Request) {
	saving := models.Saving{}

	json.NewDecoder(r.Body).Decode(&saving)

	_, err := fs.sql.Exec("INSERT INTO savings (finance_id, amount) VALUES ($1, $2)", saving.Finance, saving.Amount)
	if err != nil {
		http.Error(w, "There was an error inserting the saving", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusAccepted,
		Message: "Created the saving successfully",
		Data:    saving,
	})
}

func MakeFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{
		sql: db,
	}
}
