package models

type Wallet struct {
	Id     int64 `gorm:"column:id;primaryKey;AUTO_INCREMENT"`
	Name   string
	Amount float64
}
