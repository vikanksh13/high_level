package repo

import (
	"Interview/models"

	"github.com/jmoiron/sqlx"
)

type ITransactionRepo interface {
	IngestNewTransaction(dbConn *sqlx.DB, transaction models.Transaction) error
	GetTransactionDetails(dbConn *sqlx.DB, name string) ([]models.Transaction, error)
}

type TransactionRepo struct{}

func NewTransactionRepo() ITransactionRepo {
	return &TransactionRepo{}
}

func (repo *TransactionRepo) IngestNewTransaction(dbConn *sqlx.DB, transaction models.Transaction) error {
	query := `
		INSERT INTO transactions (wallet_id, wallet_name, before_balance, after_balance, transaction_amount, type_of_transaction)
		VALUES (:wallet_id, :wallet_name, :before_balance, :after_balance, :transaction_amount, :type_of_transaction)
	`

	_, err := dbConn.NamedExec(query, &transaction)
	if err != nil {
		return err
	}

	return nil
}

func (repo *TransactionRepo) GetTransactionDetails(dbConn *sqlx.DB, name string) ([]models.Transaction, error) {
	query := `SELECT * FROM transactions WHERE wallet_name = ?`

	rows, err := dbConn.Queryx(query, name)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction

		err := rows.StructScan(&transaction)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
