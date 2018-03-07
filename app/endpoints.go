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
var srv *http.Server = nil

func checkErrors(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		log.Print(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer checkErrors(w, r)
	var err error
	var result *string
	switch r.Method {
	case "GET":
		var d int
		if d, err = strconv.Atoi(r.URL.Path[13:]); err != nil {
		} else if result, err = GetUser(uint64(d)); err != nil {
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
			Email    string
			Password string
		}
		if err = decoder.Decode(&newUser); err != nil {
		} else if err = CreateUser(newUser.Email, newUser.Password); err != nil {
		} else if result, err = GetUser(GetUserByEmail(newUser.Email).Id); err != nil {
		}
		if err != nil {
			http.Error(w, "Error in creating new user", http.StatusInternalServerError)
			log.Print(err)
		} else {
			fmt.Fprint(w, *result)
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
		} else if result, err = GetUser(updatedUser.Id); err != nil {
		}
		if err != nil {
			http.Error(w, "Error in updating user info", http.StatusInternalServerError)
			log.Print(err)
		} else {
			fmt.Fprint(w, *result)
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer checkErrors(w, r)
	var err error
	switch r.Method {
	case "GET":
		var d int
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
		if err = decoder.Decode(&budget); err != nil {
		} else if err = CreateBudget(budget.UserId, budget.Income, budget.Rent, budget.Wealth); err != nil {
		}
		if err != nil {
			log.Print(err)
			http.Error(w, "Error in creating budget", http.StatusInternalServerError)
		} else {
			budgets := GetUserBudgets(budget.UserId)
			fmt.Fprint(w, budgets[len(budgets)-1])
		}
	case "PUT":
		var budget struct {
			Id     uint64
			Income uint64
			Rent   uint64
			Wealth int64
		}
		decoder := json.NewDecoder(r.Body)
		if err = decoder.Decode(&budget); err != nil {
		} else if err = EditBudget(budget.Id, budget.Income, budget.Rent, budget.Wealth); err != nil {
		}
		if err != nil {
			log.Print(err)
			http.Error(w, "Error in updating budget", http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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
				http.Error(w, success.Error(), http.StatusUnauthorized)
			} else {
				result, _ := json.Marshal(user)
				fmt.Fprint(w, string(result))
				UpdateUser(user.Id, user.Username, user.Email, user.EncryptedPass, user.Verified, time.Now())
			}
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer checkErrors(w, r)
	switch r.Method {
	case "GET":
		email := r.URL.Path[15:]
		if err := VerifyUser(email); err != nil {
			log.Print(err)
			http.Error(w, "Error in verifying email", http.StatusInternalServerError)
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func StartServer(port int) {
	if srv == nil {
		srv := &http.Server{Addr: fmt.Sprintf(":%v", port)}
		http.HandleFunc(baseURL+"/user/", UserHandler)
		http.HandleFunc(baseURL+"/budget/", BudgetHandler)
		http.HandleFunc(baseURL+"/verify/", VerifyUserHandler)
		http.HandleFunc(baseURL+"/login/", LoginHandler)
		err := srv.ListenAndServe()
		log.Print(err)
	} else {
		fmt.Println("Server already started")
	}
}
