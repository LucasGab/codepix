package usecase

import "github.com/LucasGab/codepix-go/domain/model"

// TransactionUseCase is ...
type TransactionUseCase struct {
	TransactioRepository model.TransactionRepositoryInterface
	PixRepository        model.PixKeyRepositoryInterface
}

// Register is ...
func (t *TransactionUseCase) Register(accountID string, amount float64, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {

	account, err := t.PixRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	err = t.TransactioRepository.Register(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Confirm is ...
func (t *TransactionUseCase) Confirm(transactionID string) (*model.Transaction, error) {

	transaction, err := t.TransactioRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Confirm()
	if err != nil {
		return nil, err
	}

	err = t.TransactioRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Complete is ...
func (t *TransactionUseCase) Complete(transactionID string) (*model.Transaction, error) {

	transaction, err := t.TransactioRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Complete()
	if err != nil {
		return nil, err
	}

	err = t.TransactioRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// Error is ...
func (t *TransactionUseCase) Error(transactionID string, reason string) (*model.Transaction, error) {

	transaction, err := t.TransactioRepository.Find(transactionID)
	if err != nil {
		return nil, err
	}

	err = transaction.Cancel(reason)
	if err != nil {
		return nil, err
	}

	err = t.TransactioRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
