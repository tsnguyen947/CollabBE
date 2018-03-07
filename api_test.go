package main

import (
	"testing"

	"github.com/tsnguyen947/CollabBE/app"
)

func init() {
	app.SwitchToTest()
	app.DB.Exec("SELECT setval('users_id_seq', 1)")
	app.DB.Exec("SELECT setval('budgets_id_seq', 1)")
	DBTeardown()
}

func TestUserGetData(t *testing.T) {
	DBSetup(t)
	userData, _ := app.GetUser(1)
	expected := "{\"Id\":1,\"Username\":\"user1\",\"Email\":\"email1@test.com\",\"EncryptedPass\":\"supersecretpass1\",\"LastAccess\":\"0001-01-01T00:00:00Z\",\"Verified\":false}"
	if *userData != expected {
		t.Errorf("User data does not match expected: %s vs %s", *userData, expected)
	}
}

func TestNilUserGetData(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic while getting nonexistent user")
		}
	}()
	app.GetUser(11)
}

func TestCreateUser(t *testing.T) {
	DBSetup(t)
	err := app.CreateUser("newguy@test.com", "hisstupidpassword")
	user := app.GetUserByEmail("newguy@test.com")
	if err != nil || user == nil {
		t.Errorf("Unexpected error when creating user: %v", err)
	}
}

func TestCreateExistingUser(t *testing.T) {
	DBSetup(t)
	err := app.CreateUser("email1@test.com", "unnecessarypass")
	if err == nil {
		t.Error("Expected error while creating user with existing email")
	}
}

func TestEditUser(t *testing.T) {
	DBSetup(t)
	app.CreateUser("newuser@test.com", "temppass")
	user := app.GetUserByEmail("newuser@test.com")
	err := app.EditUser(user.Id, "newuser1", "newuser@test.com", "temppass", "newpass", true)
	if err != nil {
		t.Error("Unexpected error while editing user")
	}
	newUser := app.GetUserByEmail("newuser@test.com")
	if *user == *newUser {
		t.Errorf("User was not updated")
	}
}

func TestEditUserWrongPass(t *testing.T) {
	DBSetup(t)
	err := app.EditUser(1, "newuser1", "newuser@test.com", "definitelythewrongpass", "newpass", true)
	if err == nil {
		t.Error("Error expected while editing user with incorrect password")
	}
}

func TestGetBudgets(t *testing.T) {
	DBSetup(t)
	budgets, err := app.GetBudgets(0)
	if err != nil {
		t.Error("Unexpected error while getting budgets")
	}
	expected := [...]string{
		"{\"Id\":0,\"UserId\":0,\"Income\":0,\"Rent\":0,\"Wealth\":0}",
		"{\"Id\":1,\"UserId\":0,\"Income\":1,\"Rent\":1,\"Wealth\":1}",
		"{\"Id\":2,\"UserId\":0,\"Income\":2,\"Rent\":2,\"Wealth\":2}",
		"{\"Id\":3,\"UserId\":0,\"Income\":3,\"Rent\":3,\"Wealth\":3}",
		"{\"Id\":4,\"UserId\":0,\"Income\":4,\"Rent\":4,\"Wealth\":4}",
	}
	for n, b := range budgets {
		if b != expected[n] {
			t.Errorf("Budget does not match expected: %s vs %s", b, expected[n])
		}
	}
}

func TestGetEmptyBudgets(t *testing.T) {
	DBSetup(t)
	budgets, err := app.GetBudgets(5)
	if err != nil {
		t.Error("Unexpected error while getting empty budget array")
	}
	if len(budgets) != 0 {
		t.Errorf("Expected empty budget array got one with %v elements", len(budgets))
	}
}

func TestGetBudgetsNilUser(t *testing.T) {
	DBSetup(t)
	defer func() {
		if err := recover(); err == nil {
			t.Error("Expected panic while getting budget from nonexistent user")
		}
	}()
	app.GetBudgets(50)
}

func TestCreateBudget(t *testing.T) {
	DBSetup(t)
	err := app.CreateBudget(5, 0, 0, 0)
	if err != nil {
		t.Error("Unexpected error received while creating budget")
	}
	budget, _ := app.GetBudgets(5)
	if len(budget) != 1 {
		t.Error("New budget not found")
	}
}

func TestCreateBudgetNilUser(t *testing.T) {
	DBSetup(t)
	err := app.CreateBudget(50, 0, 0, 0)
	if err == nil {
		t.Error("Expected error from creating budget with no user")
	}
}

func TestEditBudget(t *testing.T) {
	DBSetup(t)
	budget := app.GetBudgetById(1)
	if err := app.EditBudget(1, 50, 50, 50); err != nil {
		t.Error("Error occured while updating budget")
	}
	newBudget := app.GetBudgetById(1)
	if *budget == *newBudget {
		t.Error("Budget was not updated")
	}
}

func TestEditNilBudget(t *testing.T) {
	DBSetup(t)
	if err := app.EditBudget(50, 50, 50, 50); err == nil {
		t.Error("Expected error while editing nonexistent budget")
	}
}

func TestVerify(t *testing.T) {
	DBSetup(t)
	if err := app.CreateUser("testemail@test.com", "testpass"); err != nil {
		t.Error("Unexpected error while creating user for verification")
	}
	if err := app.VerifyUser("testemail@test.com"); err != nil {
		t.Error("Unexpected error while verifying user")
	}
	if err := app.VerifyUser("nonexistent@test.com"); err == nil {
		t.Error("Expected error while verifying nonexistent user")
	}

}

func TestLogin(t *testing.T) {
	DBSetup(t)
	app.CreateUser("testemail@test.com", "testpass")
	err := app.Login("testemail@test.com", "testpass")
	if err != nil {
		t.Error("Unexpected error from login")
	}
	err = app.Login("testemail@test.com", "incorrectpass")
	if err == nil {
		t.Error("Expected error from incorrect pass")
	}
}
