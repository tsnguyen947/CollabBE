package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var baseURL = "/api/v1"

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		d, err := strconv.Atoi(r.URL.Path[13:])
		if err != nil {
			log.Print(err)
		}
		user := GetUserByID(uint64(d))
		obj, err := json.Marshal(*user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(obj[:]))
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var data User
		err := decoder.Decode(&data)
		if err != nil {
			panic(err)
		}
		WriteUser(data.Id, data.Username, data.Rent, data.Wealth)
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func BudgetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "POST":
	default:
		http.Error(w, fmt.Sprintf("%v is not allowed at this path", r.Method), http.StatusMethodNotAllowed)
	}
}

func StartServer() {
	http.HandleFunc(baseURL+"/user/", UserHandler)
	http.ListenAndServe(":8080", nil)
}
