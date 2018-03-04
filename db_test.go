package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/tsnguyen947/CollabBE/app"
)

func init() {
	app.SwitchToTest()
	app.DB.Exec("SELECT setval('users_id_seq', 100)")
	app.DB.Exec("SELECT setval('budgets_id_seq', 100)")
	DBTeardown()
}

func checkPanics(t *testing.T) {
	if err := recover(); err != nil {
		t.Error(err)
	}
}

func DBSetup(t *testing.T) {
	defer checkPanics(t)
	DBTeardown()
	for n := uint64(0); n < 10; n++ {
		app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, false)", n, fmt.Sprintf("user%v", n), fmt.Sprintf("email%v@test.com", n), fmt.Sprintf("supersecretpass%v", n), time.Time{})
		app.DB.Exec("INSERT INTO budgets VALUES($1, $2, $3, $4, $5)", n, (n / 5), n, n, n)
	}
}

func DBTeardown() {
	app.DB.Exec("DELETE FROM budgets")
	app.DB.Exec("DELETE FROM users")
}

func TestUserGet(t *testing.T) {
	DBSetup(t)
	for n := uint64(0); n < 10; n++ {
		compare := app.User{Id: n, Username: fmt.Sprintf("user%v", n), Email: fmt.Sprintf("email%v@test.com", n), EncryptedPass: fmt.Sprintf("supersecretpass%v", n), LastAccess: time.Time{}, Verified: false}
		idUser := app.GetUserByID(n)
		emailUser := app.GetUserByEmail(fmt.Sprintf("email%v@test.com", n))
		usernameUser := app.GetUserByUsername(fmt.Sprintf("user%v", n))
		if compare != *idUser || *idUser != *emailUser || *emailUser != *usernameUser {
			t.Errorf("User does not match expected")
		}
	}
}

func TestNilUserGet(t *testing.T) {
	DBSetup(t)
	var reason []string
	if app.GetUserByID(1024) != nil {
		reason = append(reason, "ID")
	}
	if app.GetUserByEmail("nonexistent@email.com") != nil {
		reason = append(reason, "Email")
	}
	if app.GetUserByUsername("unreal_guy") != nil {
		reason = append(reason, "Username")
	}
	if len(reason) != 0 {
		t.Errorf("Issue in getting nonexistent users by %s", strings.Join(reason, " and "))
	}
}

func TestUserWrite(t *testing.T) {
	DBSetup(t)
	for n := 0; n < 10; n++ {
		app.WriteUser(fmt.Sprintf("email%v@test.com", n), fmt.Sprintf("supersecretpass%v", n))
	}
	rows, err := app.DB.Query("SELECT COUNT(*) FROM users")
	if err != nil {
		t.Error(err)
	} else {
		var n int
		rows.Next()
		rows.Scan(&n)
		if n != 20 {
			t.Errorf("Expected 20 users got %v", n)
		}
	}
}

func TestEmptyUserWrite(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected panic when creating user with empty email and pass but received none")
		}
	}()
	app.WriteUser("", "")
}

func TestUpdateUser(t *testing.T) {
	DBSetup(t)
	var n uint64 = 1
	var user, newUser app.User
	now := time.Now().UTC()
	rows, _ := app.DB.Query("SELECT * FROM users WHERE id=$1", n)
	rows.Next()
	rows.Scan(&user.Id, &user.Username, &user.Email, &user.EncryptedPass, &user.LastAccess, &user.Verified)
	compare := app.User{Id: n, Username: fmt.Sprintf("user%v", n), Email: fmt.Sprintf("email%v@test.com", n), EncryptedPass: fmt.Sprintf("supersecretpass%v", n), LastAccess: time.Time{}, Verified: false}
	if user != compare {
		t.Errorf("Error in validating user before update")
	}
	rows = nil
	app.UpdateUser(n, fmt.Sprintf("user%v", n+10), fmt.Sprintf("email%v@test.com", n+10), fmt.Sprintf("supersecretpass%v", n+10), true, now)
	rows, _ = app.DB.Query("SELECT * FROM users WHERE id=$1", n)
	rows.Next()
	rows.Scan(&newUser.Id, &newUser.Username, &newUser.Email, &newUser.EncryptedPass, &newUser.LastAccess, &newUser.Verified)
	compare = app.User{Id: n, Username: fmt.Sprintf("user%v", n+10), Email: fmt.Sprintf("email%v@test.com", n+10), EncryptedPass: fmt.Sprintf("supersecretpass%v", n+10), LastAccess: now, Verified: true}
	if newUser != compare {
		t.Errorf("Error in validating user after update")
	}
}

func TestGetUserBudgets(t *testing.T) {
	DBSetup(t)
	budgets := app.GetUserBudgets(0)
	if len(budgets) != 5 {
		t.Errorf(fmt.Sprintf("Expected 5 budgets got %v", len(budgets)))
	}
	for n, budget := range budgets {
		compare := app.Budget{Id: uint64(n), UserId: 0, Income: uint64(n), Rent: uint64(n), Wealth: uint64(n)}
		if compare != *budget {
			t.Errorf("Budget does not match expected: %v vs %v", budget, compare)
		}
	}
}

func TestGetUserBudgetsNilUser(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic while getting budgets for nonexistent user")
		}
	}()
	app.GetUserBudgets(50)
}

func TestGetEmptyUserBudgets(t *testing.T) {
	DBSetup(t)
	budgets := app.GetUserBudgets(5)
	if len(budgets) != 0 {
		t.Errorf("Expected empty budget list but got one with %v elements", len(budgets))
	}
}

func TestWriteBudget(t *testing.T) {
	DBSetup(t)
	for n := 0; n < 10; n++ {
		app.WriteBudget(7, uint64(n), uint64(n), int64(n))
	}
	rows, _ := app.DB.Query("SELECT count(*) FROM budgets")
	var n int
	rows.Next()
	rows.Scan(&n)
	if n != 20 {
		t.Errorf("Expected 20 budgets got %v", n)
	}
}

func TestWriteBudgetNilUser(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic while creating budget for nonexistent user")
		}
	}()
	app.WriteBudget(50, 10000, 10000, 10000)
}
