package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var baseURL = "/api/v1"

func checkErrors(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		log.Print(err)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	defer checkErrors(w, r)
	switch r.Method {
	case "GET":
		var d int
		var err error
		var result *string
		if d, err = strconv.Atoi(r.URL.Path[13:]); err != nil {
		} else if result, err = GetUser(d); err != nil {
		} else {
			fmt.Println(w, result)
		}
		if err != nil {
			http.Error(w, "Error in getting user data", http.StatusInternalServerError)
			log.Print(err)
		}
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var err error
		var newUser struct {
			email    string
			password string
		}
		if err = decoder.Decode(&newUser); err != nil {
		} else if err = CreateUser(newUser.email, newUser.password); err != nil {
		}
		if err != nil {
			http.Error(w, "Error in creating new user", http.StatusInternalServerError)
			log.Print(err)
		}
	case "PUT":
		var err error
		var updatedUser struct {
			Id       uint64
			Username string
			Email    string
			OldPass  string
			NewPass  string
			Verified bool
		}
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&updatedUser); err != nil {
		} else if err = EditUser(updatedUser.Id, updatedUser.Username, updatedUser.Email, updatedUser.OldPass, updatedUser.NewPass, updatedUser.Verified); err != nil {
		}
		if err != nil {
			http.Error(w, "Error in updating user info", http.StatusInternalServerError)
			log.Print(err)
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	defer checkErrors(w, r)
	switch r.Method {
	case "GET":
		var d int
		var err error
		var budgets []string
		if d, err = strconv.Atoi(r.URL.Path[15:]); err != nil {
		} else if budgets, err = GetBudgets(uint64(d)); err != nil {
		}
		if err != nil {
			log.Print(err)
			http.Error(w, "Error in getting budgets", http.StatusInternalServerError)
		} else {
			for _, str := range budgets {
				fmt.Fprint(w, str)
			}
		}
	case "POST":
		var budget struct {
			UserId uint64
			Income uint64
			Rent   uint64
			Wealth int64
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&budget); err != nil {
			log.Print(err)
			http.Error(w, "Error in creating budget", http.StatusInternalServerError)
		}
		CreateBudget(budget.UserId, budget.Income, budget.Rent, budget.Wealth)
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer checkErrors(w, r)
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var cred struct {
			Username string
			Password string
		}
		if err := decoder.Decode(&cred); err != nil {
			http.Error(w, "Error in logging in", http.StatusInternalServerError)
			log.Print(err)
		} else {
			user := GetUserByUsername(cred.Username)
			success := Login(cred.Username, cred.Password)
			if success != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				UpdateUser(user.Id, user.Username, user.Email, user.EncryptedPass, user.Verified, time.Now())
			}
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	defer checkErrors(w, r)
	switch r.Method {
	case "GET":
		email := r.URL.Path[15:]
		VerifyUser(email)
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func StartServer() {
	http.HandleFunc(baseURL+"/user/", UserHandler)
	http.HandleFunc(baseURL+"/budget/", BudgetHandler)
	http.HandleFunc(baseURL+"/verify/", VerifyUserHandler)
	http.HandleFunc(baseURL+"/login/", LoginHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
}
