package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/RinnAnd/ww-backend/config"
	"github.com/RinnAnd/ww-backend/models"
)

type PostgresDB struct {
	db *sql.DB
}

func New(cfg config.Postgres) (*PostgresDB, error) {
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %w", err)
	}

	where, err := TableMaker(db)
	if err != nil {
		fmt.Println(where)
		fmt.Println(err)
		return nil, fmt.Errorf("there was an error creating the tables: %w, at %s table", err, where)
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (p *PostgresDB) CreateUser(user *models.User) error {
	_, err := p.db.Exec(
		"INSERT INTO users (username, name, email, password) VALUES ($1, $2, $3, $4)",
		user.UserName,
		user.Name,
		user.Email,
		user.Password,
	)
	if err != nil {
		return fmt.Errorf("could not insert user: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetUserByEmail(email string) (*models.User, error) {
	user := models.User{}

	row := p.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.UserName, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("could not get user by email: %w", err)
	}

	return &user, nil
}

func (p *PostgresDB) GetAllUsers() ([]*models.User, error) {
	users := []*models.User{}

	rows, err := p.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("could not get all users: %w", err)
	}

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (p *PostgresDB) UpdatePasswordForEmail(email, password string) error {
	_, err := p.db.Exec("UPDATE users SET password = $1 WHERE email = $2", password, email)
	if err != nil {
		return fmt.Errorf("could not update password for email: %w", err)
	}

	return nil
}

func (p *PostgresDB) CreateFinance(finance *models.Finance) error {
	_, err := p.db.Exec(
		"INSERT INTO finances (user_id, month, year, salary) VALUES ($1, $2, $3, $4)",
		finance.UserID,
		finance.Month,
		finance.Year,
		finance.Salary,
	)
	if err != nil {
		return fmt.Errorf("could not insert finance: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetFinancesByUserID(userID string) ([]*models.RelFinance, error) {
	finances := []*models.RelFinance{}

	rows, err := p.db.Query("SELECT * FROM finances WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("could not get finances by user id: %w", err)
	}

	for rows.Next() {
		finance := models.RelFinance{}
		err := rows.Scan(&finance.ID, &finance.UserID, &finance.Month, &finance.Year, &finance.Salary)
		if err != nil {
			return nil, fmt.Errorf("could not scan finance: %w", err)
		}

		finances = append(finances, &finance)
	}

	return finances, nil
}

func (p *PostgresDB) CreateExpense(expense *models.Expense) error {
	_, err := p.db.Exec(
		"INSERT INTO expenses (finance_id, name, amount, category) VALUES ($1, $2, $3, $4)",
		expense.Finance,
		expense.Name,
		expense.Amount,
		expense.Category,
	)
	if err != nil {
		return fmt.Errorf("could not insert expense: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetExpensesByFinanceID(financeID string) ([]*models.Expense, error) {
	expenses := []*models.Expense{}

	rows, err := p.db.Query("SELECT * FROM expenses WHERE finance_id = $1", financeID)
	if err != nil {
		return nil, fmt.Errorf("could not get expenses by finance id: %w", err)
	}

	for rows.Next() {
		expense := models.Expense{}
		err := rows.Scan(&expense.ID, &expense.Finance, &expense.Name, &expense.Amount, &expense.Category)
		if err != nil {
			return nil, fmt.Errorf("could not scan expense: %w", err)
		}

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (p *PostgresDB) CreateSaving(saving *models.Saving) error {
	_, err := p.db.Exec(
		"INSERT INTO savings (finance_id, amount) VALUES ($1, $2)",
		saving.Finance,
		saving.Amount,
	)
	if err != nil {
		return fmt.Errorf("could not insert saving: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetSavingsByFinanceID(financeID string) ([]*models.Saving, error) {
	savings := []*models.Saving{}

	rows, err := p.db.Query("SELECT * FROM savings WHERE finance_id = $1", financeID)
	if err != nil {
		return nil, fmt.Errorf("could not get savings by finance id: %w", err)
	}

	for rows.Next() {
		saving := models.Saving{}
		err := rows.Scan(&saving.ID, &saving.Finance, &saving.Amount)
		if err != nil {
			return nil, fmt.Errorf("could not scan saving: %w", err)
		}

		savings = append(savings, &saving)
	}

	return savings, nil
}
