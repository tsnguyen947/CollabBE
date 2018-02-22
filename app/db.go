package app

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB = nil

func init() {
	var err error = nil
	db, err = sql.Open("postgres", "user=samnguyen dbname=pqgotest sslmode=require")
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserByID(id uint64) (*User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE id IS %v", id))
	if err != nil {
		return nil, err
	}
	user := new(User)
	rows.Next()
	err = rows.Scan(&user.Id, &user.Username, &user.Rent, &user.Wealth, &user.EncryptedPass)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(name string) (*User, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE username='%s'", name))
	if err != nil {
		return nil, err
	}
	user := new(User)
	rows.Next()
	err = rows.Scan(&user.Id, &user.Username, &user.Rent, &user.Wealth, &user.EncryptedPass)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func WriteUser(id uint64, username string, rent uint64, wealth uint64, encryptedPass string) {
	_, err := db.Exec("INSERT INTO users VALUES ($1, $2, $3, $4, $5)", id, username, rent, wealth, encryptedPass)
	if err != nil {
		log.Print(err)
	}
}

func GetUserBudgets(userID uint64) ([]*Budget, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM budgets WHERE userId IS %v", userID))
	if err != nil {
		return nil, err
	}
	budgets := make([]*Budget, 5)
	for rows.Next() {
		budget := new(Budget)
		err = rows.Scan(&budget.Id, &budget.UserId, &budget.Other)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}
