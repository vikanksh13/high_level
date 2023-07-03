package service

import (
	"Interview/models"
	"Interview/repo"
	"errors"

	"github.com/jmoiron/sqlx"
)

type WalletOperationsService struct {
	dbConn          *sqlx.DB
	WalletRepo      repo.IWalletRepo
	TranscationRepo repo.ITransactionRepo
}

type IWalletOperationsService interface {
	CreateWallet(walletValue models.Wallet) error
	UpdateWalletAmount(amount float64, name string) error
	GetHistory(name string) ([]models.Transaction, error)
}

func NewWalletOperationsService(dbConn *sqlx.DB, WalletRepo repo.IWalletRepo, TranscationRepo repo.ITransactionRepo) IWalletOperationsService {
	return &WalletOperationsService{
		WalletRepo:      WalletRepo,
		TranscationRepo: TranscationRepo,
	}
}

func (svc *WalletOperationsService) CreateWallet(walletValue models.Wallet) error {

	//calling repo to check this name is already in the wallet or not
	modelVal, err := svc.WalletRepo.GetByName(svc.dbConn, walletValue.Name)
	if err != nil {
		return err
	}
	if modelVal != nil {
		return errors.New("this name is already in the wallet")
	}

	//insert a new wallet in the repo
	err = svc.WalletRepo.IngestNewUser(svc.dbConn, walletValue)
	if err != nil {
		return err
	}

	return nil
}

func (svc *WalletOperationsService) UpdateWalletAmount(amount float64, name string) error {
	modelVal, err := svc.WalletRepo.GetByName(svc.dbConn, name)
	if err != nil {
		return err
	}

	err = svc.WalletRepo.UpdateAmount(svc.dbConn, modelVal.Amount+amount)
	// If the concurrent update are happening what we can do is if the amount is greater than 0 that means credited, than we can perform update concurrently but
	// if it is less than 0 that means debited then first we can get the amount then acquire lock(mutex) onto it then perform the operation after that unlock it.
	if err != nil {
		return err
	}

	var transactionType string
	if amount > 0 {
		transactionType = "credit"
	} else {
		transactionType = "debit"
	}

	transcationVal := models.Transaction{
		BeforeBalance:     modelVal.Amount,
		AfterBalance:      modelVal.Amount + amount,
		TransactionAmount: amount,
		TypeOfTransaction: transactionType,
	}

	err = svc.TranscationRepo.IngestNewTransaction(svc.dbConn, transcationVal)
	if err != nil {
		return err
	}

	return nil
}

func (svc *WalletOperationsService) GetHistory(name string) ([]models.Transaction, error) {

	transcationVal, err := svc.TranscationRepo.GetTransactionDetails(svc.dbConn, name)
	if err != nil {
		return nil, err
	}

	return transcationVal, nil
}
