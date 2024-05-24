package models

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	Email    string
	Password string
}

type Session struct {
	ID       string
	UserName string
	Name     string
	Email    string
}

type Friendship struct {
	ID      string `json:"id"`
	UserID1 string `json:"user_id1"`
	UserID2 string `json:"user_id2"`
}

type Finance struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
	Salary int    `json:"salary"`
}

type FullFinance struct {
	ID       string
	UserID   string
	Month    int
	Year     int
	Salary   int
	Expenses []Expense
}

type Expense struct {
	ID       string `json:"id"`
	Finance  string `json:"finance_id"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	Category string `json:"category"`
}

type Saving struct {
	ID      string `json:"id"`
	Finance string `json:"finance_id"`
	Amount  int    `json:"amount"`
}
