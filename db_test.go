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
	app.DB.Exec("SELECT setval('users_id_seq', 1)")
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
}

func DBTeardown() {
	app.DB.Exec("DELETE FROM users")
	app.DB.Exec("DELETE FROM budgets")
}

func TestUserGet(t *testing.T) {
	DBSetup(t)
	for n := uint64(0); n < 10; n++ {
		app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, true)", n, fmt.Sprintf("user%v", n), fmt.Sprintf("email%v@test.com", n), fmt.Sprintf("supersecretpass%v", n), time.Time{})
		compare := app.User{Id: n, Username: fmt.Sprintf("user%v", n), Email: fmt.Sprintf("email%v@test.com", n), EncryptedPass: fmt.Sprintf("supersecretpass%v", n), LastAccess: time.Time{}, Verified: true}
		idUser := app.GetUserByID(n)
		emailUser := app.GetUserByEmail(fmt.Sprintf("email%v@test.com", n))
		usernameUser := app.GetUserByUsername(fmt.Sprintf("user%v", n))
		if compare != *idUser || *idUser != *emailUser || *emailUser != *usernameUser {
			t.Errorf("User does not match expected")
		}
	}
	DBTeardown()
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
	DBTeardown()
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
		if n != 10 {
			t.Errorf("Expected 10 got %v", n)
		}
	}
	DBTeardown()
}

func TestEmptyUserWrite(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Expected panic when creating user with empty email and pass but received none")
		}
	}()
	app.WriteUser("", "")
	DBTeardown()
}

func TestUpdateUser(t *testing.T) {
	DBSetup(t)
	var n uint64 = 1
	var user, newUser app.User
	now := time.Now().UTC()
	app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, false)", n, fmt.Sprintf("user%v", n), fmt.Sprintf("email%v@test.com", n), fmt.Sprintf("supersecretpass%v", n), time.Time{})
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
	DBTeardown()
}

func TestGetBudgets(t *testing.T) {
	DBSetup(t)
	app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, true)", 1, "user", "email@test.com", "supersecretpass", time.Time{})
	for n := uint64(0); n < 10; n++ {
		app.DB.Exec("INSERT INTO budgets VALUES($1, $2, $3, $4, $5)", n, 1, n, n, n)
	}
	budgets := app.GetUserBudgets(1)
	if len(budgets) != 10 {
		t.Errorf(fmt.Sprintf("Expected 10 rows got %v", len(budgets)))
	}
	for n, budget := range budgets {
		compare := app.Budget{Id: uint64(n), UserId: 1, Income: uint64(n), Rent: uint64(n), Wealth: uint64(n)}
		if compare != *budget {
			t.Errorf("Budget does not match expected: %v vs %v", budget, compare)
		}
	}
	DBTeardown()
}

func TestGetBudgetsNilUser(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic but got none")
		}
	}()
	for n := uint64(0); n < 10; n++ {
		app.DB.Exec("INSERT INTO budgets VALUES($1, $2, $3, $4, $5)", n, 1, n, n, n)
	}
	app.GetUserBudgets(1)
	DBTeardown()
}

func TestGetEmptyBudgets(t *testing.T) {
	DBSetup(t)
	app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, true)", 1, "user", "email@test.com", "supersecretpass", time.Time{})
	budgets := app.GetUserBudgets(1)
	if len(budgets) != 0 {
		t.Errorf("Expected empty budget list but got one with %v elements", len(budgets))
	}
	DBTeardown()
}

func TestCreateBudget(t *testing.T) {
	DBSetup(t)
	app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, true)", 1, "user", "email@test.com", "supersecretpass", time.Time{})
	for n := 0; n < 10; n++ {
		app.WriteBudget(1, uint64(n), uint64(n), int64(n))
	}
	rows, _ := app.DB.Query("SELECT count(*) FROM budgets")
	var n int
	rows.Next()
	rows.Scan(&n)
	if n != 10 {
		t.Errorf("Expected 10 got %v", n)
	}
	DBTeardown()
}

func TestCreateBudgetNilUser(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic but got none")
		}
	}()
	app.WriteBudget(1, 10000, 10000, 10000)
	DBTeardown()
}
