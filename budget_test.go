package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tsnguyen947/CollabBE/app"
)

func init() {
	app.SwitchToTest()
}

func DBSetup() {
	for n := 0; n < 10; n++ {
		app.DB.Exec("INSERT INTO users VALUES($1, $2, $3, $4, $5, true)", n, fmt.Sprintf("user%v", n), fmt.Sprintf("email%v@test.com", n), fmt.Sprintf("supersecretpass%v", n), time.Time{})
	}
}

func DBTeardown() {
	app.DB.Exec("DELETE FROM users")
}

func TestUserGet(t *testing.T) {
	DBSetup()
	var n uint64
	for n = 0; n < 10; n++ {
		compare := app.User{Id: n, Username: fmt.Sprintf("user%v", n), Email: fmt.Sprintf("email%v@test.com", n), EncryptedPass: fmt.Sprintf("supersecretpass%v", n), LastAccess: time.Time{}, Verified: true}
		user := app.GetUserByID(n)
		user.LastAccess = compare.LastAccess
		if compare != *user {
			t.Errorf("User received does not match expected")
		}
	}
	DBTeardown()
}
