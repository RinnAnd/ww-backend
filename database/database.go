package database

import (
	"fmt"

	"github.com/RinnAnd/ww-backend/config"
	"github.com/RinnAnd/ww-backend/database/postgres"
	"github.com/RinnAnd/ww-backend/models"
)

type Database interface {
	CreateUser(*models.User) error
	GetUserByEmail(string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdatePasswordForEmail(string, string) error
	CreateFinance(*models.Finance) error
	GetFinancesByUserID(string) ([]*models.RelFinance, error)
	CreateExpense(*models.Expense) error
	GetExpensesByFinanceID(string) ([]*models.Expense, error)
	CreateSaving(*models.Saving) error
	GetSavingsByFinanceID(string) ([]*models.Saving, error)
}

func New(cfg config.Database) (Database, error) {
	switch cfg.Driver {
	case config.PostgresDBDriver:
		return postgres.New(cfg.Postgres)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
