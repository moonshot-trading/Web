package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type GetQuote struct {
	StockSymbol    string
	UserId         string
	TransactionNum int
}

type GetUser struct {
	UserId         string
	TransactionNum int
}

type Dumplog struct {
	Filename       string
	TransactionNum int
}

type StockValue struct {
	UserId         string
	StockSymbol    string
	Amount         int
	TransactionNum int
}

type AddFunds struct {
	UserId         string
	Amount         int
	TransactionNum int
}

type ReturnQuote struct {
	Stock          string
	Price          int
	TransactionNum int
}

type webConfig struct {
	ts string
}

var stockName = regexp.MustCompile("([a-zA-Z][a-zA-Z][a-zA-Z])|([a-zA-Z])")
var config = webConfig{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST")
	} else {
		return "localhost"
	}
}()}

//ioutil.Discard
//os.Stdout
var Log = log.New(ioutil.Discard, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)

func checkValidStockName(s string) bool {
	m := stockName.FindStringSubmatch(s)
	if m == nil {
		Log.Println("bad stock name")
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
	response.TransactionNum = d.TransactionNum
	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func respondQuoteRequests(w http.ResponseWriter, d GetQuote) {

	response := ReturnQuote{}
	response.Stock = d.StockSymbol
	response.Price = 123
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	Log.Println("Quote", d)

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
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func addFundsHandler(w http.ResponseWriter, r *http.Request) {

	var d AddFunds
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Log.Println("Add funds", d)
	if r.Method == "POST" {
		if d.UserId == "" {
			http.Error(w, "No User", http.StatusInternalServerError)
			return
		}
	}

	//hit TS to check that user exists
	//maybe send back if they are new or existing

	sendToTServer(d, "add")

	response := AddFunds{}
	response.UserId = d.UserId
	response.Amount = d.Amount
	response.TransactionNum = d.TransactionNum

	_, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func buyHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("Buy", d)
	//check user account
	sendToTServer(d, "buy")

	respondStockValueRequests(w, d)
}

func commitBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}
	sendToTServer(d, "confirmBuy")
	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId
	response.TransactionNum = d.TransactionNum
	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func cancelBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("cancel buy", d)

	sendToTServer(d, "cancelBuy")
	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func sellHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		Log.Println("sell error", d)
		return
	}

	Log.Println("sell ", d)
	sendToTServer(d, "sell")

	//check user account

	respondStockValueRequests(w, d)
}

func commitSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("commit sell", d)
	sendToTServer(d, "confirmSell")

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func cancelSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("cancel sell", d)
	sendToTServer(d, "cancelSell")

	//hit TS to check that user exists
	//confirm buy etc
	//not sure what sending back w yet

	response := GetUser{}
	response.UserId = d.UserId
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func setBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("Set Buy", d)
	sendToTServer(d, "setBuy")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	Log.Println("Cancel Set Buy", d)
	sendToTServer(d, "cancelSetBuy")

	respondQuoteRequests(w, d)
}

func buyTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("buy trigger", d)
	sendToTServer(d, "setBuyTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func setSellHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("Set Sell", d)
	sendToTServer(d, "setSell")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	Log.Println("Cancel Set Sell", d)
	sendToTServer(d, "cancelSetSell")

	respondQuoteRequests(w, d)
}

func sellTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("sell trigger", d)
	sendToTServer(d, "setSellTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func displaySummaryHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("diaply summary", d)
	sendToTServer(d, "displaySummary")
	//the response type will be different later it will be a bunch of data
	response := GetUser{}
	response.UserId = d.UserId
	response.TransactionNum = d.TransactionNum

	_, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func dumplogHandler(w http.ResponseWriter, r *http.Request) {

	var d Dumplog

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Log.Println("dump", d)
	sendToTServer(d, "dumpLog")

	//the response type will be different later it will be a bunch of data
	response := Dumplog{}
	response.Filename = d.Filename
	response.TransactionNum = d.TransactionNum

	_, err = json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func failOnError(err error, msg string) {
	if err != nil {
		Log.Printf("%s: %s", msg, err)
		panic(err)
	}
}

func sendToTServer(r interface{}, s string) {
	jsonValue, _ := json.Marshal(r)
	resp, err := http.Post("http://"+config.ts+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
	failOnError(err, "Error sending request tp TS")
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
}

func runningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}
	return false
}

func main() {
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../frontend"))))
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
	http.HandleFunc("/Dumplog", dumplogHandler)
	http.HandleFunc("/DisplaySummary", displaySummaryHandler)

	http.ListenAndServe(":8080", nil)
}
