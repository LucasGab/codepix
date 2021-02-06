package factory

import (
	"github.com/LucasGab/codepix-go/application/usecase"
	"github.com/LucasGab/codepix-go/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

// TransactionUseCaseFactory is a ...
func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	transactionRepository := repository.TransactionRepositoryDb{Db: database}

	transactionUseCase := usecase.TransactionUseCase{
		TransactioRepository: transactionRepository,
		PixRepository:        pixRepository,
	}

	return transactionUseCase
}
