package model

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Transaction is a ...
type Transaction struct {
	ID           string  `json:"id" validate:"required,uuid4"`
	AccountID    string  `json:"accountId" validate:"required,uuid4"`
	Amount       float64 `json:"amount" validate:"required,numeric"`
	PixKeyTo     string  `json:"pixKeyTo" validate:"required"`
	PixKeyKindTo string  `json:"pixKeyKindTo" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	Status       string  `json:"status" validate:"-"`
	Error        string  `json:"error"`
}

// isValid is a ...
func (t *Transaction) isValid() error {
	v := validator.New()
	err := v.Struct(t)
	if err != nil {
		fmt.Errorf("Error during Transaction validation: %s", err.Error())
		return err
	}
	return nil
}

// ParseJSON is a ...
func (t *Transaction) ParseJSON(data []byte) error {
	err := json.Unmarshal(data, t)

	if err != nil {
		return err
	}

	err = t.isValid()
	if err != nil {
		return err
	}

	return nil
}

// ToJSON is a ...
func (t *Transaction) ToJSON() ([]byte, error) {
	err := t.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(t)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// NewTransaction is a ...
func NewTransaction() *Transaction {
	return &Transaction{}
}
