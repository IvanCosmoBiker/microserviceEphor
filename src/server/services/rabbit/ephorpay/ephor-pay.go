package ephorpay

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
	connectionPostgresql "ephorservices/src/pkg/connectDb"
    counter "ephorservices/src/pkg/counter"
	configEphor "ephorservices/src/configs"
	ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
)
	const (
		VendState_Session	 		= 1 //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
		VendState_Approving   		= 2 //[14] Продукт выбран. Ожидание оплаты.
		VendState_Vending	 		= 3 //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
		VendState_VendOk	     	= 4 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
		VendState_VendError     	= 5 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
	)

	const (
		VendError_VendFailed         = "769" //769 Ошибка выдачи продукта
		VendError_SessionCancelled   = "770" //770
		VendError_SessionTimeout     = "771" //771
		VendError_WrongProduct       = "772" //772
		VendError_VendCancelled      = "773" //773
		VendError_ApprovingTimeout   = "774" //774
	)

	const (
		TransactionState_Idle 				= 0; // Transaction Idle
		TransactionState_MoneyHoldStart 	= 1; // создали транзакцию банка
		TransactionState_MoneyHoldWait 		= 2; // ожидает ответ от банка
		TransactionState_VendSessionStart 	= 3; // PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
		TransactionState_VendSession	 	= 4; //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
		TransactionState_VendApproving   	= 5; //[14] Продукт выбран. Ожидание оплаты.
		TransactionState_Vending	 		= 6; //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
		TransactionState_MoneyDebitStart	= 8;
		TransactionState_MoneyDebitWait		= 9;
		TransactionState_MoneyDebitOk		= 10;
		TransactionState_VendOk				= 11; // приготовил продукт
		TransactionState_Error 				= 120;
	) 

type Request struct {
	Tid     int
	St 		int
	D       string
	Err 	int
	A 		int
	Wid 	int
	Sum		int
}

func initPayCron(forever chan bool) {
	stringQueue := cfg.Services.EphorPay.NameQueue
	go func() {
		select {
		case <-forever:
		    ConnectionRabbit.CloseChannel(stringQueue)
			if counterGo.N == 0 {
				return 
			}
		}
	}()
	msg, _ := ConnectionRabbit.RabbitMQConsume(stringQueue)
	counterGo.Add()
	for d := range msg {
		req := Request{}
		log.Printf("\n [x] %s", d.Body)
		dataLog := fmt.Sprintf("%s", d.Body)
		err2 := json.Unmarshal(d.Body, &req)

		if err2 != nil {
			errData, _ := fmt.Println(err2)
			log.Println(errData)
			ConnDb.AddLog(dataLog, "EphorPay", fmt.Sprintf("%s", err2),"ephorPayError")
			Transactions.Send(req.Tid,d.Body)
			continue
		}
		Transactions.Send(req.Tid,d.Body)
		//log.Printf("\n [x] %v",result)
		//log.Printf("\n [x] %v",req.Tid)
		//log.Printf("\n [x] %+v",Transactions.GetTransactions())
		ConnDb.AddLog(dataLog,"EphorPay", " ",req.D)
		//checkStatusTransaction(&req)
	}
	counterGo.Sub()
	return 
}

var cfg *configEphor.Config
var ConnectionRabbit *ConnectionRabbitMQ.ChannelMQ
var ConnDb connectionPostgresql.DatabaseInstance
var counterGo counter.Counter
var forever = make(chan bool)
var Transactions *transactionDispetcher.TransactionDispetcher

func Start(conf *configEphor.Config,Rabbit *ConnectionRabbitMQ.ChannelMQ,Db connectionPostgresql.DatabaseInstance,transactions *transactionDispetcher.TransactionDispetcher) {
	fmt.Println("Start EphorPay...")
	ConnectionRabbit = Rabbit
	ConnDb = Db
	cfg = conf
	Transactions = transactions
	start(forever)
}

func ReconnectQueue() {
	initPayCron(forever)
}

func start(forever chan bool) {
	initPayCron(forever)
}

func stop(forever chan bool) {
	ConnectionRabbit.CloseConnectRabbit()
	if counterGo.N == 0 {
		ConnDb.CloseConnectionDb()
	} else {
		go func() {
			select {
			case <-time.After(10 * time.Second):
				forever <- true
			}
		}()
		ConnDb.CloseConnectionDb()
	}
	os.Exit(3)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getDescriptionCode(code int) string {
	var stringCode string = ""
	 switch code {
		 case VendState_Session:
		 stringCode = `Оплата успешна, ожидание нажатия пользователем кнопки на ТА`
		 return stringCode 
		 fallthrough
		 case VendState_Approving:
		 stringCode = `Продукт выбран. Ожидание оплаты.`
		 return stringCode
		 fallthrough
		 case VendState_Vending:
		 stringCode = `Оплата успешна, ТА готовит продукт`
		 return stringCode
		 fallthrough
		 case VendState_VendOk:
		 stringCode = `Оплата успешна, ТА приготовил продукт`
		 return stringCode
		 fallthrough
		 case VendState_VendError:
		 stringCode = `Ошибка`
		 return stringCode
	 }
	 return stringCode
}

func getStatusServer(status int) int {
	 switch status {
		 case VendState_Session:
		 return TransactionState_VendSession
		 fallthrough
		 case VendState_Approving:
		 return TransactionState_VendApproving
		 fallthrough
		 case VendState_Vending:
		 return TransactionState_Vending
		 fallthrough
		 case VendState_VendOk:
		 return TransactionState_VendOk
		 fallthrough
		 case VendState_VendError:
		 return TransactionState_Error
	 }
	 return TransactionState_Error
}

func getDescriptionErr(err interface {}) string{
	errCode := fmt.Sprintf("%v",err)
	stringErr := ""
	 switch errCode {
		 case VendError_VendFailed:
		 stringErr = `Ошибка выдачи товара на автомате`
		 return stringErr
		 fallthrough
		 case VendError_SessionCancelled:
		 stringErr = `Продажа отменена автоматом`
		 return stringErr
		 fallthrough
		 case VendError_SessionTimeout:
		 stringErr = `Время ожидание выбора товара на автомате истекло`
		 return stringErr
		 fallthrough
		 case VendError_WrongProduct:
		 stringErr = `Выбранный на автомате товар не совпадает с оплаченым`
		 return stringErr
		 fallthrough
		 case VendError_VendCancelled:
		 stringErr = `Выдача товара отменена автоматом`
		 return stringErr
		 fallthrough
		 case VendError_ApprovingTimeout:
		 stringErr = `Время ожидание оплаты истекло`
		 return stringErr
	 }
	 return stringErr
}

func getTransactionMap(req *Request) (map[string]interface{}, map[string]interface{}) {
	Where := make(map[string]interface{})
	parametrs := make(map[string]interface{})
	Where["id"] = req.Tid
	parametrs["status"] = getStatusServer(req.St)
	parametrs["error"] = req.Err
	return Where,parametrs
}

func checkStatusTransaction(req *Request){
	Where,parametrs := getTransactionMap(req)
	errorCode := parametrs["error"]
	stringErr := ""
	if req.Err != 0 {
		stringErr = getDescriptionErr(errorCode)
	}
	if stringErr == "" {
		stringCode := getDescriptionCode(req.St)
		log.Printf("%s",stringCode)
		parametrs["error"] = stringCode
	}else {
		parametrs["error"] = stringErr
	}
	ConnDb.Set("transaction", parametrs, Where)
}
