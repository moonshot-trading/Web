package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type GetQuote struct {
	StockSymbol string
	UserId      string
}

type GetUser struct {
	UserId string
}

type StockValue struct {
	UserId      string
	StockSymbol string
	Amount      int
}

type AddFunds struct {
	UserId string
	Amount int
}

type ReturnQuote struct {
	Stock string
	Price int
}

var stockName = regexp.MustCompile("([a-zA-Z][a-zA-Z][a-zA-Z])|([a-zA-Z])")

func checkValidStockName(s string) bool {
	m := stockName.FindStringSubmatch(s)
	if m == nil {
		fmt.Println("bad stock name")
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
		if d.UserId == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
		return true
	} else {
		return false
	}
}

func verifyStockValueRequests(w http.ResponseWriter, r *http.Request, d *StockValue) bool {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Println("yeet")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if r.Method == "POST" {

		pass := checkValidStockName(d.StockSymbol)
		if pass == false {
			http.Error(w, "Wrong stock symbol", http.StatusInternalServerError)
			return false
		}

		if d.Amount == 0 {
			http.Error(w, "No money", http.StatusInternalServerError)
			return false
		}
		if d.UserId == "" {
			http.Error(w, "No User", http.StatusInternalServerError)
			return false
		}
		return true
	} else {
		return false
	}
}

func verifyQuoteRequests(w http.ResponseWriter, r *http.Request, d *GetQuote) bool {
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	if r.Method == "POST" {

		pass := checkValidStockName(d.StockSymbol)
		if pass == false {
			http.Error(w, "Wrong stock symbol", http.StatusInternalServerError)
			return false
		}

		if d.UserId == "" {
			http.Error(w, "Empty User", http.StatusInternalServerError)
			return false
		}

		return true
	} else {
		return false
		//not post
	}
}

//add some error passing from TS
func respondStockValueRequests(w http.ResponseWriter, d StockValue) {

	response := StockValue{}
	response.StockSymbol = d.StockSymbol
	response.Amount = 123 //this amount will be rounded and returned in TS
	response.UserId = d.UserId

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func respondQuoteRequests(w http.ResponseWriter, d GetQuote) {

	response := ReturnQuote{}
	response.Stock = d.StockSymbol
	response.Price = 123

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	fmt.Println("Quote", d)

	sendToTServer(d, "quote")

	respondQuoteRequests(w, d)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	//hit TS to check that user exists
	//maybe send back if they are new or existing

	response := GetUser{}
	response.UserId = d.UserId

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func addFundsHandler(w http.ResponseWriter, r *http.Request) {

	var d AddFunds
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Add funds", d)
	if r.Method == "POST" {
		if d.UserId == "" {
			http.Error(w, "No User", http.StatusInternalServerError)
			return
		}
		if d.Amount < 1 {
			http.Error(w, "Negative funds cannot be added", http.StatusInternalServerError)
			return
		}
	}

	//hit TS to check that user exists
	//maybe send back if they are new or existing

	sendToTServer(d, "add")

	response := AddFunds{}
	response.UserId = d.UserId
	response.Amount = d.Amount
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
	fmt.Println("Buy", d)
	//check user account
	sendToTServer(d, "buy")

	respondStockValueRequests(w, d)
}

func commitBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}
	fmt.Println("commit buy", d)
	sendToTServer(d, "confirmBuy")
	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId

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

	fmt.Println("cancel buy", d)

	sendToTServer(d, "cancelBuy")
	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId

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
		fmt.Println("sell error", d)
		return
	}

	fmt.Println("sell ", d)
	sendToTServer(d, "sell")

	//check user account

	respondStockValueRequests(w, d)
}

func commitSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	fmt.Println("commit sell", d)
	sendToTServer(d, "confirmSell")

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId

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

	fmt.Println("cancel sell", d)
	sendToTServer(d, "cancelSell")

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func setBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	fmt.Println("Set Buy", d)
	sendToTServer(d, "setBuy")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	fmt.Println("Cancel Set Buy", d)
	sendToTServer(d, "cancelSetBuy")

	respondQuoteRequests(w, d)
}

func buyTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	fmt.Println("buy trigger", d)
	sendToTServer(d, "setBuyTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func setSellHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	fmt.Println("Set Sell", d)
	sendToTServer(d, "setSell")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	fmt.Println("Cancel Set Sell", d)
	sendToTServer(d, "cancelSetSell")

	respondQuoteRequests(w, d)
}

func sellTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	fmt.Println("sell trigger", d)
	sendToTServer(d, "setSellTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func displaySummaryHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	fmt.Println("diaply summary", d)
	//the response type will be different later it will be a bunch of data
	response := GetUser{}
	response.UserId = d.UserId

	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseJson)
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(err)
	}
}

func sendToTServer(r interface{}, s string) *http.Response {
	jsonValue, _ := json.Marshal(r)
	resp, err := http.Post("http://localhost:44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
	failOnError(err, "Error sending request tp TS")
	return resp
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../frontend"))))
	http.HandleFunc("/GetQuote", quoteHandler)
	http.HandleFunc("/AddUser", addUserHandler)
	http.HandleFunc("/AddFunds", addFundsHandler)
	http.HandleFunc("/BuyStock", buyHandler)
	http.HandleFunc("/CommitBuy", commitBuyHandler)
	http.HandleFunc("/CancelBuy", cancelBuyHandler)
	http.HandleFunc("/SellStock", sellHandler)
	http.HandleFunc("/CancelSell", cancelSellHandler)
	http.HandleFunc("/CommitSell", commitSellHandler)
	http.HandleFunc("/SetBuyAmount", setBuyHandler)
	http.HandleFunc("/SetBuyTrigger", buyTriggerHandler)
	http.HandleFunc("/CancelSetBuy", cancelSetBuyHandler)
	http.HandleFunc("/SetSellAmount", setSellHandler)
	http.HandleFunc("/SetSellTrigger", sellTriggerHandler)
	http.HandleFunc("/CancelSetSell", cancelSetSellHandler)
	http.HandleFunc("/DisplaySummary", displaySummaryHandler)

	http.ListenAndServe(":8080", nil)
}
