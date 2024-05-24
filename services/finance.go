package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RinnAnd/ww-backend/models"
	"github.com/RinnAnd/ww-backend/utils"
	"github.com/gorilla/mux"
)

type FinanceService struct {
	sql *sql.DB
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

func (fs *FinanceService) GetUserFinances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]

	finances := []models.FullFinance{}

	rows, err := fs.sql.Query(`
        SELECT f.*, e.* 
        FROM finances f
        LEFT JOIN expenses e ON f.id = e.finance_id
        WHERE f.user_id = $1`, userId)
	if err != nil {
		http.Error(w, "There was an error fetching the user's finances", http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		finance := models.FullFinance{}
		expense := models.Expense{}
		err := rows.Scan(&finance.ID, &finance.UserID, &finance.Month, &finance.Year, &finance.Salary, &expense.ID, &expense.Finance, &expense.Name, &expense.Amount, &expense.Category)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "There was an error fetching finance", http.StatusBadRequest)
			return
		}

		finance.Expenses = append(finance.Expenses, expense)
		finances = append(finances, finance)
	}

	message := fmt.Sprintf("Here are the finances for the user %s", userId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response{
		Status:  http.StatusOK,
		Message: message,
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

func MakeFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{
		sql: db,
	}
}
