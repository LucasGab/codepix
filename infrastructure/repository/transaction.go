package repository

import (
	"fmt"

	"github.com/LucasGab/codepix-go/domain/model"
	"gorm.io/gorm"
)

// type TransactionRepositoryInterface interface {
// 	Register(transaction *Transaction) error
// 	Save(transaction *Transaction) error
// 	Find(id string) (*Transaction, error)
// }

// TransactionRepositoryDb is ...
type TransactionRepositoryDb struct {
	Db *gorm.DB
}

// Register is ...
func (t TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	err := t.Db.Create(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Save is ...
func (t TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	err := t.Db.Save(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

// Find is ...
func (t TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	t.Db.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("no transaction was found")
	}
	return &transaction, nil
}
