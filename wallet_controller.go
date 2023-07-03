package main

import (
	"Interview/models"
	"Interview/repo"
	"Interview/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type WalletController struct {
	dbConn          *sqlx.DB
	WalletRepo      repo.IWalletRepo
	TranscationRepo repo.ITransactionRepo
	WalletService   service.IWalletOperationsService
}

func NewWalletController(apiMux *mux.Router, dbConn *sqlx.DB, WalletRepo repo.IWalletRepo,
	TranscationRepo repo.ITransactionRepo, WalletService service.IWalletOperationsService) *WalletController {
	return &WalletController{
		dbConn:          dbConn,
		WalletRepo:      WalletRepo,
		TranscationRepo: TranscationRepo,
		WalletService:   WalletService,
	}

}

func (ctrl *WalletController) CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	var wallet models.Wallet

	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = ctrl.WalletService.CreateWallet(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctrl *WalletController) UpdateWalletAmountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	amount := vars["amount"]
	name := vars["name"]

	// Convert the amount string to float64
	amountValue, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	err = ctrl.WalletService.UpdateWalletAmount(amountValue, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ctrl *WalletController) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	history, err := ctrl.WalletService.GetHistory(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(history)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
