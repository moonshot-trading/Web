package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type GetQuote struct {
	Stock string
	User  string
}

type GetUser struct {
	User string
}

type StockValue struct {
	User   string
	Stock  string
	Amount int
}

type ReturnQuote struct {
	Stock string
	Price int
}

var stockName = regexp.MustCompile("([a-zA-Z][a-zA-Z][a-zA-Z])")

func checkValidStockName(s string) bool {
	m := stockName.FindStringSubmatch(s)
	if m == nil {
		return false
	}
	return true
}

func verifyGetUser(w http.ResponseWriter, r *http.Request, d *GetUser) bool {

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if r.Method == "POST" {
		if d.User == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
	}
	return true
}

func verifyStockValueRequests(w http.ResponseWriter, r *http.Request, d *StockValue) bool {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println("yeet")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if r.Method == "POST" {

		pass := checkValidStockName(d.Stock)
		if pass == false {
			http.Error(w, "Wrong stock symbol", http.StatusInternalServerError)
			return false
		}

		if d.Amount == 0 {
			http.Error(w, "No money", http.StatusInternalServerError)
			return false
		}
		if d.User == "" {
			http.Error(w, "No User", http.StatusInternalServerError)
			return false
		}
	}
	return true
}

//add some error passing from TS
func respondStockValueRequests(w http.ResponseWriter, d StockValue) {

	response := StockValue{}
	response.Stock = d.Stock
	response.Amount = 123 //this amount will be rounded and returned in TS
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if r.Method == "POST" {

		pass := checkValidStockName(d.Stock)
		if pass == false {
			http.Error(w, "Wrong stock symbol", http.StatusInternalServerError)
			return
		}

		if d.User == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ReturnQuote{}
		response.Stock = d.Stock
		response.Price = 123

		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(responseJson)
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//maybe send back if they are new or existing

	response := GetUser{}
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	//check user account

	respondStockValueRequests(w, d)
}

func commitBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func cancelBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func sellHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}

	//check user account

	respondStockValueRequests(w, d)
}

func commitSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func cancelSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.User = d.User

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func buyTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}

	//check user account

	respondStockValueRequests(w, d)

}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../frontend"))))
	http.HandleFunc("/GetQuote", quoteHandler)
	http.HandleFunc("/AddUser", addUserHandler)
	http.HandleFunc("/BuyStock", buyHandler)
	http.HandleFunc("/CommitBuy", commitBuyHandler)
	http.HandleFunc("/CancelBuy", cancelBuyHandler)
	http.HandleFunc("/CancelSell", cancelSellHandler)
	http.HandleFunc("/CommitSell", commitSellHandler)
	http.HandleFunc("/SellStock", buyHandler)
	http.HandleFunc("/BuyStockTrigger", buyTriggerHandler)

	http.ListenAndServe(":8080", nil)
}
