package kafka

import (
	"fmt"
	"os"

	"github.com/LucasGab/codepix-go/application/factory"
	appmodel "github.com/LucasGab/codepix-go/application/model"
	"github.com/LucasGab/codepix-go/application/usecase"
	"github.com/LucasGab/codepix-go/domain/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
)

// Processor is a ...
type Processor struct {
	Database     *gorm.DB
	Producer     *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

// NewProcessor is a ...
func NewProcessor(database *gorm.DB, producer *ckafka.Producer, deliveryChan chan ckafka.Event) *Processor {
	return &Processor{
		Database:     database,
		Producer:     producer,
		DeliveryChan: deliveryChan,
	}
}

// Consume is a ...
func (p *Processor) Consume() {
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
		"group.id":          os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset": "earliest",
	}
	c, err := ckafka.NewConsumer(configMap)

	if err != nil {
		panic(err)
	}

	topics := []string{os.Getenv("kafkaTransactionTopic"), os.Getenv("kafkaTransactionConfirmationTopic")}
	c.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Println(string(msg.Value))
			p.processMessage(msg)
		}
	}
}

// processMessage is a ...
func (p *Processor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"

	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		p.processTransaction(msg)
	case transactionConfirmationTopic:
		p.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (p *Processor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(p.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyKindTo,
		transaction.Description,
		transaction.ID,
	)

	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending
	transactionJSON, err := transaction.ToJSON()

	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), topic, p.Producer, p.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

// processTransactionConfirmation is a ...
func (p *Processor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJSON(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(p.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = p.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
		return nil
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// confirmTransaction is a ...
func (p *Processor) confirmTransaction(transaction *appmodel.Transaction, transactionUseCase usecase.TransactionUseCase) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)

	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	transactionJSON, err := transaction.ToJSON()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJSON), topic, p.Producer, p.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}
