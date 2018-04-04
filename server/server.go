package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type webConfig2 struct {
	ts2 string
}

type webConfig3 struct {
	ts3 string
}

type webConfig4 struct {
	ts4 string
}

type webConfig5 struct {
	ts5 string
}

var stockName = regexp.MustCompile("([a-zA-Z][a-zA-Z][a-zA-Z])|([a-zA-Z])")

var semaphoreChan = make(chan struct{}, 80)

var config = webConfig{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST")
	} else {
		return "localhost"
	}
}()}

var config2 = webConfig2{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST_2")
	} else {
		return "localhost"
	}
}()}

var config3 = webConfig3{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST_3")
	} else {
		return "localhost"
	}
}()}

var config4 = webConfig4{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST_4")
	} else {
		return "localhost"
	}
}()}

var config5 = webConfig5{func() string {
	if runningInDocker() {
		return os.Getenv("TX_SERVER_HOST_5")
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

	sendToTServer(d, d.UserId, "quote")

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

	sendToTServer(d, d.UserId, "add")

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
	sendToTServer(d, d.UserId, "buy")

	respondStockValueRequests(w, d)
}

func commitBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}
	sendToTServer(d, d.UserId, "confirmBuy")
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

	sendToTServer(d, d.UserId, "cancelBuy")
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
	sendToTServer(d, d.UserId, "sell")

	//check user account

	respondStockValueRequests(w, d)
}

func commitSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("commit sell", d)
	sendToTServer(d, d.UserId, "confirmSell")

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
	sendToTServer(d, d.UserId, "cancelSell")

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
	sendToTServer(d, d.UserId, "setBuy")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetBuyHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	Log.Println("Cancel Set Buy", d)
	sendToTServer(d, d.UserId, "cancelSetBuy")

	respondQuoteRequests(w, d)
}

func buyTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("buy trigger", d)
	sendToTServer(d, d.UserId, "setBuyTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func setSellHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("Set Sell", d)
	sendToTServer(d, d.UserId, "setSell")

	//check user account

	respondStockValueRequests(w, d)
}

func cancelSetSellHandler(w http.ResponseWriter, r *http.Request) {

	var d GetQuote

	if !verifyQuoteRequests(w, r, &d) {
		return
	}
	Log.Println("Cancel Set Sell", d)
	sendToTServer(d, d.UserId, "cancelSetSell")

	respondQuoteRequests(w, d)
}

func sellTriggerHandler(w http.ResponseWriter, r *http.Request) {

	var d StockValue
	if !verifyStockValueRequests(w, r, &d) {
		return
	}
	Log.Println("sell trigger", d)
	sendToTServer(d, d.UserId, "setSellTrigger")

	//check user account
	respondStockValueRequests(w, d)
}

func displaySummaryHandler(w http.ResponseWriter, r *http.Request) {

	var d GetUser
	if !verifyGetUser(w, r, &d) {
		return
	}

	Log.Println("diaply summary", d)
	sendToTServer(d, d.UserId, "displaySummary")
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
	sendToTServer(d, "", "dumpLog")

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
		fmt.Printf("%s: %s\n", msg, err)
		panic(err)
	}
}

func sendToTServer(r interface{}, user string, s string) {
	server := config.ts
	jsonValue, _ := json.Marshal(r)

	semaphoreChan <- struct{}{}
	go func() {
		if user == "" {
			fmt.Println("dumpppp", r)
			resp, err := http.Post("http://"+config.ts+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
			failOnError(err, "Error sending request to TS")
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()

			resp2, err2 := http.Post("http://"+config2.ts2+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
			failOnError(err2, "Error sending request to TS")
			io.Copy(ioutil.Discard, resp2.Body)
			resp2.Body.Close()

			resp3, err3 := http.Post("http://"+config3.ts3+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
			failOnError(err3, "Error sending request to TS")
			io.Copy(ioutil.Discard, resp3.Body)
			resp3.Body.Close()

			resp4, err4 := http.Post("http://"+config4.ts4+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
			failOnError(err4, "Error sending request to TS")
			io.Copy(ioutil.Discard, resp4.Body)
			resp4.Body.Close()

			resp5, err5 := http.Post("http://"+config5.ts5+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
			failOnError(err5, "Error sending request to TS")
			io.Copy(ioutil.Discard, resp5.Body)
			resp5.Body.Close()

			<-semaphoreChan
			return

		} else if len(user) > 0 && user[0] < 'C' { //1/3 of the way
			server = config.ts
		} else if len(user) > 0 && user[0] < 'P' {
			server = config2.ts2
		} else if len(user) > 0 && user[0] < 'c' {
			server = config3.ts3
		} else if len(user) > 0 && user[0] < 'o' {
			server = config4.ts4
		} else {
			server = config5.ts5
		}

		resp, err := http.Post("http://"+server+":44416/"+s, "application/json", bytes.NewBuffer(jsonValue))
		failOnError(err, "Error sending request to TS")
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()

		<-semaphoreChan
	}()
}

func runningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	if err == nil {
		return true
	}
	return false
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/go/src/github.com/moonshot-trading/Web/frontend"))))
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
