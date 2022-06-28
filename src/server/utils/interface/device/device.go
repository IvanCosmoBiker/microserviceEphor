package device

import (
	connectionPostgresql "ephorservices/src/pkg/connectDb"
	configEphor "ephorservices/src/configs"
	ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    transactionStruct "ephorservices/src/data/transaction"
    interfaseBank "ephorservices/src/server/utils/interface/payment"
	requestApi "ephorservices/src/data/requestApi"
    transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
)

var (
    TypeCoffee = 0
	TypeSnack = 1
	TypeHoreca = 2
	TypeSodaWater = 3
	TypeMechanical = 4
	TypeComb = 5
	TypeMicromarket = 6
	TypeCooler = 7
)

type Device interface {
    InitDeviceData(transactionStruct.Transaction)
    SendMessage(trasactionStruct transactionStruct.Transaction, 
	conn connectionPostgresql.DatabaseInstance, 
	rebbit *ConnectionRabbitMQ.ChannelMQ, 
	conf *configEphor.Config, 
	req requestApi.Request, 
	bank interfaseBank.Bank,
	resultBank map[string]interface{},
	transactionDispetcer transactionDispetcher.TransactionDispetcher,
	channel chan []byte) map[string]interface{}
}