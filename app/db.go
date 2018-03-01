package app

import (
	"database/sql"
	"fmt"
	"log"
	"time"

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

func GetUserByAttribute(att string, val string) *User {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM users WHERE %s=%s", att, val))
	if err != nil {
		panic(err)
	}
	user := new(User)
	exists := rows.Next()
	if exists {
		if err = rows.Scan(&user.Id, &user.Username, &user.Email, &user.EncryptedPass, &user.LastAccess, &user.Verified); err != nil {
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

func WriteUser(id uint64, email string, encryptedPass string) {
	if _, err := db.Exec("INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6)", id, email, email, encryptedPass, time.Now(), false); err != nil {
		panic(err)
	}
}

func UpdateUser(id uint64, username string, email string, encryptedPass string, verified bool, access time.Time) {
	if access.IsZero() {
		if _, err := db.Exec("UPDATE users SET username=$1, email=$2, encryptedPass=$3, verified=$4 WHERE id=$5", username, email, encryptedPass, verified, id); err != nil {
			panic(err)
		}
	} else {
		if _, err := db.Exec("UPDATE users SET username=$1, email=$2, encryptedPass=$3, lastAccess=$4, verified=$5 WHERE id=$6", username, email, encryptedPass, access, verified, id); err != nil {
			panic(err)
		}
	}
}

func GetUserBudgets(userID uint64) []*Budget {
	var rows *sql.Rows
	var err error
	if rows, err = db.Query(fmt.Sprintf("SELECT * FROM budgets WHERE userId IS %v", userID)); err != nil {
		panic(err)
	}
	budgets := make([]*Budget, 5)
	for rows.Next() {
		budget := new(Budget)
		if err := rows.Scan(&budget.Id, &budget.UserId, &budget.Income, &budget.Rent, &budget.Wealth); err != nil {
			panic(err)
		}
		budgets = append(budgets, budget)
	}
	return budgets
}

func WriteBudget(userID uint64, income uint64, rent uint64, wealth int64) {
	if _, err := db.Exec("INSERT INTO budgets VALUES ($1, $2, $3, $4, $5)", 0, userID, income, rent, wealth); err != nil {
		panic(err)
	}
}
