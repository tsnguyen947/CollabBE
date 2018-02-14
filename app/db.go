package app

import (
	"bytes"
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

func GeneralGet(table string, args map[string]string) interface{} {
	conditions := bytes.NewBufferString("")
	for k, v := range args {
		conditions.WriteString(" " + k)
		conditions.WriteString("=" + v + " AND")
	}
	conditions.Truncate(conditions.Len() - 3)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s where%s", table, conditions.String()))
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	var result interface{}
	switch table {
	case "users":
		result = populateUsers(rows)
	case "budgets":
	}
	return result
}

func populateUsers(rows *sql.Rows) []*User {
	result := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.Id, &user.Username, &user.Rent, &user.Wealth)
		if err != nil {
			log.Print(err)
		}
		result = append(result, user)
	}
	return result
}

func GetUserByID(id uint64) *User {
	return GetUser(map[string]string{"id": fmt.Sprintf("%v", id)})[0]
}

func GetUser(args map[string]string) []*User {
	return GeneralGet("users", args).([]*User)
}

func WriteUser(id uint64, username string, rent uint64, wealth uint64) {
	_, err := db.Exec("INSERT INTO users VALUES ($1, $2, $3, $4)", id, username, rent, wealth)
	if err != nil {
		log.Print(err)
	}
}

// func GetBudgets(userID uint64)

func WriteBook(isbn string, title string, author string, price float32) {
	db.Exec("INSERT INTO books VALUES ($1, $2, $3, $4)", isbn, title, author, price)
}
