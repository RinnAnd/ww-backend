package models

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Finance struct {
	ID       string   `json:"id"`
	UserID   string   `json:"user_id"`
	Month    int      `json:"month"`
	Year     int      `json:"year"`
	Salary   int      `json:"salary"`
	Expenses []string `json:"expenses_ids"`
	Saving   string   `json:"savings_id"`
}

type Expense struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Ammount  int    `json:"ammount"`
	Category string `json:"category"`
}

type Saving struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
}
