package repository

import (
	"database/sql"

	"gorm.io/gorm"
)

type TxBeginner interface {
	Begin(opts ...*sql.TxOptions) (Transactioner, error)
}

type txBeginner struct {
	db *gorm.DB
}

func NewTxBeginner(db *gorm.DB) TxBeginner {
	return &txBeginner{db: db}
}

func (t *txBeginner) Begin(opts ...*sql.TxOptions) (Transactioner, error) {
	tx := t.db.Begin(opts...)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &transactioner{db: tx}, nil

}
