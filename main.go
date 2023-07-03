package main

import (
	"Interview/config"
	"Interview/repo"
	"Interview/service"
	"fmt"
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

// Setting up the payment wallet
// -> Wallet name(should be unique), default balance should be 0 and should be generating a unique wallet id.
// -> trsaction amount > 0 => credit and trsaction amount < 0 => debit  there should be a wallet id
// -> history of transcation like a passbook.
// 		balanceBefore	balanceAfter	trancsationAmount    Type of transaction

func main() {
	SetupDatabase()

	dbConn := config.GetDBConn()
	router := mux.NewRouter()
	walletRepo := repo.NewWalletRepo()
	transactionRepo := repo.NewTransactionRepo()
	walletSeervice := service.NewWalletOperationsService(dbConn, walletRepo, transactionRepo)
	walletController := NewWalletController(router, dbConn, walletRepo, transactionRepo, walletSeervice)

	router.HandleFunc("/wallet", walletController.CreateWalletHandler).Methods("POST")
	router.HandleFunc("/wallet/{name}/{amount}", walletController.UpdateWalletAmountHandler).Methods("POST")
	router.HandleFunc("/wallet/{name}/history", walletController.GetHistoryHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func SetupDatabase() error {

	err := config.OpenDatabase()
	if err != nil {
		fmt.Println("Could not open database")
	}
	return nil
}
