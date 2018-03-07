package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB = nil

func init() {
	var err error = nil
	if DB, err = sql.Open("postgres", "user=samnguyen dbname=budget sslmode=require"); err != nil {
		log.Fatal(err)
	}
}

func SwitchToTest() {
	var err error = nil
	if DB, err = sql.Open("postgres", "user=samnguyen dbname=budgettest sslmode=require"); err != nil {
		log.Fatal(err)
	}
}

func GetUserByAttribute(att string, val string) *User {
	var rows *sql.Rows
	var err error
	if rows, err = DB.Query(fmt.Sprintf("SELECT * FROM users WHERE %s=%s", att, val)); err != nil {
		panic(err)
	}
	user := new(User)
	exists := rows.Next()
	if exists {
		if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.EncryptedPass, &user.LastAccess, &user.Verified); err != nil {
			panic(err)
		}
		return user
	} else {
		return nil
	}
}

func GetUserByID(id uint64) *User {
	return GetUserByAttribute("id", fmt.Sprintf("%v", id))
}

func GetUserByEmail(email string) *User {
	return GetUserByAttribute("email", fmt.Sprintf("'%s'", email))
}

func GetUserByUsername(username string) *User {
	return GetUserByAttribute("username", fmt.Sprintf("'%s'", username))
}

func WriteUser(email string, encryptedPass string) {
	if email == "" || encryptedPass == "" {
		panic(errors.New("Cannot have an empty email or password"))
	}
	if _, err := DB.Exec("INSERT INTO users VALUES (DEFAULT, $1, $2, $3, $4, $5)", email, email, encryptedPass, time.Now(), false); err != nil {
		panic(err)
	}
}

func UpdateUser(id uint64, username string, email string, encryptedPass string, verified bool, access time.Time) {
	if access.IsZero() {
		if _, err := DB.Exec("UPDATE users SET username=$1, email=$2, encryptedPass=$3, verified=$4 WHERE id=$5", username, email, encryptedPass, verified, id); err != nil {
			panic(err)
		}
	} else {
		if _, err := DB.Exec("UPDATE users SET username=$1, email=$2, encryptedPass=$3, lastAccess=$4, verified=$5 WHERE id=$6", username, email, encryptedPass, access, verified, id); err != nil {
			panic(err)
		}
	}
}

func GetUserBudgets(userID uint64) []*Budget {
	var rows *sql.Rows
	var err error
	if user := GetUserByID(userID); user == nil {
		panic(errors.New("Attempted to access budgets of a nonexistent user"))
	}
	if rows, err = DB.Query(fmt.Sprintf("SELECT * FROM budgets WHERE userId=%v", userID)); err != nil {
		panic(err)
	}
	var budgets []*Budget
	for rows.Next() {
		budget := new(Budget)
		if err := rows.Scan(&budget.Id, &budget.UserId, &budget.Income, &budget.Rent, &budget.Wealth); err != nil {
			panic(err)
		}
		budgets = append(budgets, budget)
	}
	return budgets
}

func GetBudgetById(budgetId uint64) *Budget {
	var rows *sql.Rows
	var err error
	if rows, err = DB.Query(fmt.Sprintf("SELECT * FROM budgets WHERE id=%v", budgetId)); err != nil {
		panic(err)
	}
	budget := new(Budget)
	exists := rows.Next()
	if exists {
		if err := rows.Scan(&budget.Id, &budget.UserId, &budget.Income, &budget.Rent, &budget.Wealth); err != nil {
			panic(err)
		}
		return budget
	} else {
		return nil
	}
}

func UpdateBudget(budgetId uint64, income uint64, rent uint64, wealth int64) {
	if _, err := DB.Exec("UPDATE budgets SET income=$1, rent=$2, wealth=$3 WHERE id=$4", income, rent, wealth, budgetId); err != nil {
		panic(err)
	}
}

func WriteBudget(userID uint64, income uint64, rent uint64, wealth int64) {
	if user := GetUserByID(userID); user == nil {
		panic(errors.New("Attempted to create a budget for a nonexistent user"))
	}
	if _, err := DB.Exec("INSERT INTO budgets VALUES (DEFAULT, $1, $2, $3, $4)", userID, income, rent, wealth); err != nil {
		panic(err)
	}
}
