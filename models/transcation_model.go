package models

type Transaction struct {
	TransactionId     int64 `gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	WalletId          int64
	WalletName        string //FK
	BeforeBalance     float64
	AfterBalance      float64
	TransactionAmount float64
	TypeOfTransaction string
}
