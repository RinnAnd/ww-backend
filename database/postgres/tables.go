package postgres

import "database/sql"

func TableMaker(db *sql.DB) (string, error) {
	var err error
	_, err = db.Exec(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS users (
		    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				username TEXT NOT NULL,
		    name TEXT NOT NULL,
		    email TEXT NOT NULL UNIQUE,
		    password TEXT NOT NULL
		)`)
	if err != nil {
		return "users", err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS friendships (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        user_id1 UUID,
        user_id2 UUID,
        FOREIGN KEY (user_id1) REFERENCES users (id),
        FOREIGN KEY (user_id2) REFERENCES users (id)
    )`)
	if err != nil {
		return "friendships", err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS finances (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				user_id UUID,
				month INTEGER,
				year INTEGER,
				salary INTEGER,
				FOREIGN KEY (user_id) REFERENCES users (id)
		)`)
	if err != nil {
		return "finances", err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS savings (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			finance_id UUID,
			FOREIGN KEY (finance_id) REFERENCES finances (id),
			amount INTEGER
	)`)
	if err != nil {
		return "savings", err
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS expenses (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        finance_id UUID,
        FOREIGN KEY (finance_id) REFERENCES finances (id),
        name TEXT NOT NULL,
        amount INTEGER NOT NULL,
        category TEXT NOT NULL
    )`)
	if err != nil {
		return "expenses", err
	}

	return "all good", nil
}
