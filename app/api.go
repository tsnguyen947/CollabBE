package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var baseURL = "/api/v1"

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		d, err := strconv.Atoi(r.URL.Path[13:])
		if err != nil {
			log.Print(err)
		}
		user, err := GetUserByID(uint64(d))
		obj, err := json.Marshal(*user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(obj[:]))
	case "PUT":
		decoder := json.NewDecoder(r.Body)
		var data struct {
			Id       uint64
			Username string
			Rent     uint64
			Wealth   uint64
			Password string
		}
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		user, err := GetUserByUsername(data.Username)
		if user == nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
			if err != nil {
				panic(err)
			}
			WriteUser(data.Id, data.Username, data.Rent, data.Wealth, string(hashedPassword))
		} else {
			http.Error(w, fmt.Sprintf("Username %s is already taken", data.Username), http.StatusBadRequest)
		}
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		d, err := strconv.Atoi(r.URL.Path[15:])
		if err != nil {
			log.Print(err)
		}
		budgets, err := GetUserBudgets(uint64(d))
		for _, b := range budgets {
			obj, err := json.Marshal(b)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, fmt.Sprintln(string(obj[:])))
		}
	case "PUT":
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var cred struct {
			Username string
			Password string
		}
		err := decoder.Decode(&cred)
		if err != nil {
			log.Print(err)
			http.Error(w, fmt.Sprint("Parameters accepted are 'username' and 'password'"), http.StatusBadRequest)
		}
		user, err := GetUserByUsername(cred.Username)
		if err == nil {
			err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPass), []byte(cred.Password))
		}
		if err != nil {
			log.Print(err)
			http.Error(w, fmt.Sprint("Username or password is invalid"), http.StatusBadRequest)
		}
	}
}

func StartServer() {
	http.HandleFunc(baseURL+"/user/", UserHandler)
	http.HandleFunc(baseURL+"/budget/", BudgetHandler)
	http.HandleFunc(baseURL+"/login/", LoginHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
}
