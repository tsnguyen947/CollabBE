package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUser(userId int) (*string, error) {
	user := GetUserByID(uint64(userId))
	var obj []byte
	var err error
	if obj, err = json.Marshal(*user); err != nil {
		return nil, err
	}
	result := string(obj[:])
	return &result, nil
}

func CreateUser(email string, password string) error {
	user := GetUserByEmail(email)
	if user == nil {
		var hashedPassword []byte
		var err error
		if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
			return err
		}
		WriteUser(email, string(hashedPassword))
		return nil
	} else {
		return errors.New(fmt.Sprintf("User with email %s already exists", email))
	}
}

func EditUser(id uint64, username string, email string, oldPass string, newPass string, verified bool) error {
	user := GetUserByID(id)
	fmt.Println(user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPass), []byte(oldPass)); err != nil {
		return err
	}
	if newPass == "" {
		UpdateUser(id, username, email, user.EncryptedPass, verified, time.Time{})
	} else {
		var hashedPassword []byte
		var err error
		if hashedPassword, err = bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost); err != nil {
			return err
		}
		UpdateUser(id, username, email, string(hashedPassword), verified, time.Time{})
	}
	return nil
}

func GetBudgets(id uint64) ([]string, error) {
	budgets := GetUserBudgets(id)
	var result []string
	for _, b := range budgets {
		var obj []byte
		var err error
		if obj, err = json.Marshal(b); err != nil {
			return nil, err
		}
		result = append(result, string(obj))
	}
	return result, nil
}

func CreateBudget(userID uint64, income uint64, rent uint64, wealth int64) error {
	user := GetUserByID(userID)
	if user != nil {
		WriteBudget(userID, income, rent, wealth)
		return nil
	} else {
		return errors.New("Cannot create a budget for a nonexistent user")
	}
}

func Login(username string, password string) error {
	user := GetUserByUsername(username)
	var status error = nil
	if user != nil || !user.Verified {
		if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPass), []byte(password)); err != nil {
			status = errors.New("Username or password is invalid")
		}
	} else {
		status = errors.New("Username or password is invalid, or email is not verified")
	}
	return status
}

func VerifyUser(email string) error {
	user := GetUserByEmail(email)
	if user != nil {
		UpdateUser(user.Id, user.Username, email, user.EncryptedPass, true, time.Now())
		return nil
	} else {
		return errors.New("Email given is not a valid email")
	}
}
