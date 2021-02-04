package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

// Bank is ...
type Bank struct {
	Base     `valid:"required"`
	Code     string     `json:"code" gorm:"type:varchar(20);notnull" valid:"notnull"`
	Name     string     `json:"name" gorm:"type:varchar(255);notnull" valid:"notnull"`
	Accounts []*Account `gorm:"ForeignKey:BankID" valid:"-"`
}

// Is valid is ...
func (bank *Bank) isValid() error {
	_, err := govalidator.ValidateStruct(bank)
	if err != nil {
		return err
	}
	return nil
}

// NewBank is ...
func NewBank(code string, name string) (*Bank, error) {
	bank := Bank{
		Code: code,
		Name: name,
	}

	bank.ID = uuid.NewV4().String()
	bank.CreatedAt = time.Now()

	err := bank.isValid()
	if err != nil {
		return nil, err
	}

	return &bank, nil
}
