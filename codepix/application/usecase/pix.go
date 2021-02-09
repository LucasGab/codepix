package usecase

import "github.com/LucasGab/codepix-go/domain/model"

// PixUseCase is ...
type PixUseCase struct {
	PixKeyRepository model.PixKeyRepositoryInterface
}

// RegisterKey is ...
func (p *PixUseCase) RegisterKey(key string, kind string, accountID string) (*model.PixKey, error) {
	account, err := p.PixKeyRepository.FindAccount(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, key, account)
	if err != nil {
		return nil, err
	}

	_, err = p.PixKeyRepository.RegisterKey(pixKey)
	if err != nil {
		return nil, err
	}

	return pixKey, nil
}

// FindKey is ...
func (p *PixUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	pixKey, err := p.PixKeyRepository.FindKeyByKind(key, kind)
	if err != nil {
		return nil, err
	}
	return pixKey, nil
}
