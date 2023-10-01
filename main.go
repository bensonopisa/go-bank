package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	handler    *mux.Router
}
type apiFunc func(w http.ResponseWriter, r *http.Request) error

var routes Routes
var accounts Accounts

func init() {
	accounts = []Account{}

	routes = []Route{
		{
			Path:        "/account/create",
			Method:      http.MethodPost,
			HandlerFunc: convertToHttpHandler(handleCreateAccount),
			Description: "create new account",
		},
		{
			Path:        "/account/{id}",
			Method:      http.MethodGet,
			HandlerFunc: convertToHttpHandler(handleGetAccountById),
			Description: "fetch account by id",
		},
		{
			Path:        "/account/update",
			Method:      http.MethodPut,
			HandlerFunc: convertToHttpHandler(handleUpdateAccount),
			Description: "Update account information",
		},
		{
			Path:        "/account/delete/{id}",
			Method:      http.MethodDelete,
			HandlerFunc: convertToHttpHandler(handleAccountDelete),
			Description: "Handle deleting of an account",
		},
		{
			Path:        "/accounts/",
			Method:      http.MethodGet,
			HandlerFunc: convertToHttpHandler(handleGetAllAccounts),
			Description: "Retrieve all the accounts",
		},
	}
}

func NewServer(addr string, router *mux.Router) *ApiServer {
	return &ApiServer{
		listenAddr: addr,
		handler:    router,
	}
}

func (s *ApiServer) Run() {
	log.Println("Starting api server at address:", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, s.handler))
}

// api handler to http.HandlerFunc converter

func convertToHttpHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// might need proper error handling here...
			log.Println("an error was encountered: ", err)
		}
	}
}

// defining the handlers

func handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, accounts)
}
func handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	account, err := serializeBody(r)
	if err != nil {
		return err
	}
	accounts = append(accounts, *NewAccount(account.Name))
	return writeJSON(w, accounts)
}

func handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	account, err := serializeBody(r)
	if err != nil {
		return err
	}

	for _, v := range accounts {
		if v.ID == account.ID {
			log.Println("Account name ", account.Name)
			v.Name = account.Name;
			return writeJSON(w, &v)
		}
	}
	return nil
}
func handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	if id, ok := mux.Vars(r)["id"]; ok {
		identifier, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		for _, v := range accounts {
			if v.ID == identifier {
				return writeJSON(w, v)
			}
		}
		return errors.New("no account was found matching id=" + id)
	}
	return errors.New("please include a valid request id")
}
func handleAccountDelete(w http.ResponseWriter, r *http.Request) error {
	if id, ok := mux.Vars(r)["id"]; ok {
		log.Println("Deleting account with id", id)
		identifier, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		for idx, v := range accounts {
			if v.ID == identifier {
				accounts = append(accounts[:idx], accounts[idx+1:]...)
				return writeJSON(w, accounts)
			}
		}
		// will handle this when we get to the database
		return nil
	}
	return fmt.Errorf("please provide the account id as path variable")
}

func writeJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func serializeBody(r *http.Request) (account Account, err error) {
	err = json.NewDecoder(r.Body).Decode(&account)
	return
}

func main() {
	router := mux.NewRouter()

	for _, v := range routes {
		router.Path(v.Path).Methods(v.Method).HandlerFunc(v.HandlerFunc)
	}

	s := NewServer(":3000", router)
	s.Run()
}
