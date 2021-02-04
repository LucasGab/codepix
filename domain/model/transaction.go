package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	// TransactionPending is ...
	TransactionPending string = "pending"
	// TransactionCompleted is ...
	TransactionCompleted string = "completed"
	// TransactionError is ...
	TransactionError string = "error"
	// TransactionConfirmed is ...
	TransactionConfirmed string = "confirmed"
)

// TransactionRepositoryInterface is ....
type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

// Transactions is ...
type Transactions struct {
	Transactions []Transaction
}

// Transaction is ...
type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	AccountFromID     string   `gorm:"column:account_from_id;type:uuid;not null" valid:"notnull"`
	Amount            float64  `json:"amount" gorm:"type:float;not null" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyIDTo        string   `gorm:"column:pix_key_id_to;type:uuid;not null" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20);not null" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255);not null" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" gorm:"type:varchar(255);not null" valid:"-"`
}

// Is Valid is ...
func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if transaction.Status != TransactionCompleted && transaction.Status != TransactionPending &&
		transaction.Status != TransactionConfirmed && transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}
	return nil
}

// Complete is ...
func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}

// Cancel is ...
func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.UpdatedAt = time.Now()
	transaction.CancelDescription = description
	err := transaction.isValid()
	return err
}

// Confirm is ...
func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdatedAt = time.Now()
	err := transaction.isValid()
	return err
}

// NewTransaction is ...
func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
