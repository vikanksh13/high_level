package repo

import (
	"Interview/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IWalletRepo interface {
	IngestNewUser(dbConn *sqlx.DB, user models.Wallet) error
	UpdateAmount(dbConn *sqlx.DB, amount float64) error
	GetByName(dbConn *sqlx.DB, name string) (*models.Wallet, error)
}

type WalletRepo struct{}

func NewWalletRepo() IWalletRepo {
	return &WalletRepo{}
}

func (r *WalletRepo) IngestNewUser(dbConn *sqlx.DB, user models.Wallet) error {

	query := `
		INSERT INTO wallets (name, amount)
		VALUES (:name, :amount)
	`

	// Execute the query with the user data
	_, err := dbConn.NamedExec(query, &user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *WalletRepo) UpdateAmount(dbConn *sqlx.DB, amount float64) error {

	query := `
		UPDATE wallets
		SET amount = :amount
	`

	// Execute the query with the new amount
	_, err := dbConn.NamedExec(query, map[string]interface{}{
		"amount": amount,
	})
	if err != nil {
		return err
	}

	return nil
}

func (repo *WalletRepo) GetByName(dbConn *sqlx.DB, name string) (*models.Wallet, error) {
	query := `SELECT * FROM wallets WHERE name = ?`

	// Execute the query and retrieve the wallet
	var wallet models.Wallet
	err := dbConn.Get(&wallet, query, name)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
