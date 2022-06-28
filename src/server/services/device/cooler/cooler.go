package cooler 

import (
    "encoding/json"
    "fmt"
    "time"
    "log"
    "strconv"
    transactionStruct "ephorservices/src/data/transaction"
    interfaceDevice "ephorservices/src/server/utils/interface/device"
    ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    interfaseBank "ephorservices/src/server/utils/interface/payment"
    requestApi "ephorservices/src/data/requestApi"
    transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
    configEphor "ephorservices/src/configs"
    transactionProduct "ephorservices/src/server/model/schema/main/transactionproduct"
    automatEvent "ephorservices/src/server/model/schema/account/automatevent"
    automat "ephorservices/src/server/model/schema/account/automat"
)

const (
    IceboxStatus_Drink = 1 // выдача напитков
    IceboxStatus_Icebox = 2 // дверь открыта
    IceboxStatus_End = 3 // выдача завершена
)

type Cooler struct {
    Name string
    Status int
}

type NewCoolerStruct struct {
    Cooler
}

type RequestRabbit struct {
	Tid     int
	St 		int
	D       string
	Err 	int
	A 		int
	Wid 	int
	Sum		int
}

func GetDateTime(stringTime string,seconds int) string {
    t,_ := time.Parse("2006-01-02 15:04:05",stringTime)
    t2 := t.Add(-180 * time.Minute)
    newTime := t2.Add(time.Duration(seconds) * time.Second)
    resultTime := newTime.Format("2006-01-02 15:04:05")
    return resultTime
}

func GateDateTimeMoscow(stringTime string,seconds int) string {
     t,_ := time.Parse("2006-01-02 15:04:05",stringTime)
    newTime1 :=  t.Add(-3 * time.Hour)
    newTime := newTime1.Add(time.Duration(seconds) * time.Second)
    resultTime := newTime.Format("2006-01-02 15:04:05")
    return resultTime
}

func SetDataTransaction(parametrs,where map[string]interface{})  {
    connectDb.Set("transaction",parametrs,where)
    connectDb.SetLog(fmt.Sprintf("%q",parametrs["error"]))
}

func SendMessageToDevice(jsonData []byte,bank interfaseBank.Bank,Request requestApi.Request) bool {
     tid := fmt.Sprintf("%v",trasactionData.Tid)
     err := Rabbit.PublishMessage(jsonData,fmt.Sprintf("ephor.1.dev.%v",Request.Imei))
     if err != nil {
         bank.ReturnMoney(tid,Request)
         return false
     }
     return true
}

func WaitResponseRabbit(channel chan []byte,Bank interfaseBank.Bank,Request requestApi.Request) RequestRabbit {
    timer := time.NewTimer(time.Duration(Config.RabbitMq.ExecuteTimeSeconds) * time.Second)
    tid := fmt.Sprintf("%v",trasactionData.Tid)
    log.Printf("\n [x] %s","Timer")
    requestRabbit := RequestRabbit{}
    select {
        case <-timer.C:{
            where["id"] = Request.IdTransaction
            parametrs["ps_desc"] = "Время истекло"
            parametrs["error"] = "Время истекло"
            parametrs["status"] = transactionStruct.TransactionState_Error
            SetDataTransaction(parametrs,where)
            Bank.ReturnMoney(tid,Request)
            fmt.Println("Timer is end")
            requestRabbit = RequestRabbit{
                St: transactionStruct.TransactionState_ErrorTimeOut,
                Err: transactionStruct.TransactionState_ErrorTimeOut,
            }
            timer.Stop()
        }
        case result := <-channel: {
            timer.Stop()
            log.Printf("\n [x] %s","Timer ok")
            log.Println("Timer ok")
            json.Unmarshal(result, &requestRabbit)
        }
    }
    return requestRabbit
}

func AddDataTransactionProducts(transaction transactionStruct.Transaction){
    params := make(map[string]interface{})
    for _, product := range transaction.Products {
        params["transaction_id"] = transaction.Tid
        params["name"] = product["name"]
        params["select_id"] = product["select_id"]
        params["ware_id"] = product["ware_id"]
        params["value"] = product["price"]
        params["quantity"] = product["quantity"]
        TransactionProductStore.AddByParams(params)
    }
}

func AddDataAutomatEvent(transaction transactionStruct.Transaction){
    params := make(map[string]interface{})
    addSeconds := 1
    for _, product := range transaction.Products {
        quantity := product["quantity"].(float64)
        quantityInt := int(quantity)
        fmt.Println(quantity)
        for i := 0; i < quantityInt; i++ {
            date := GateDateTimeMoscow(transaction.Date,addSeconds)
            params["account_id"] = transaction.AccountId
            params["automat_id"] = transaction.AutomatId
            params["type"] = 3
            params["date"] = date
            params["category"] = 1
            params["credit"] = product["price"]
            params["name"] = product["name"]
            params["select_id"] = product["select_id"]
            params["ware_id"] = product["ware_id"]
            params["value"] = product["price"]
            params["status"] = 0
            if transaction.PointId != 0 {
                params["point_id"] = transaction.PointId
            }else {
                params["point_id"] = nil
            }
            AutomatEventStore.AddByParams(params)
            addSeconds +=1
        }
    }
}

func GetAutomat(transaction transactionStruct.Transaction) ([]map[string]interface{}){
    options := make(map[string]interface{})
    options["id"] = transaction.AutomatId
    options["account_id"] = transaction.AccountId
    automats := AutomatStore.GetWithOptions(options)
    return automats
}

func UpdateAutomat(transaction transactionStruct.Transaction){
    log.Println("UpdateAutomat")
    automat := GetAutomat(transaction)
    if len(automat) == 0 {
        return
    }
    date := time.Now()
    params := make(map[string]interface{})
    products := len(transaction.Products)
    for _, entry := range automat {
        cashlessNum,_ := entry["now_cashless_num"].(int64)
        cashlessVal,_ := entry["now_cashless_val"].(int64)
        params["account_id"] = transaction.AccountId
        params["id"] = transaction.AutomatId
        params["now_cashless_val"] = cashlessVal + int64(transaction.Sum)
        params["now_cashless_num"] = cashlessNum + int64(products)
        params["last_sale"] = date.Format("2006-01-02 15:04:05")
        params["last_login"] = date.Format("2006-01-02 15:04:05")
        AutomatStore.SetByParams(params)
    }  
}

 // пока не используется.
func (cl Cooler) StartPostpaid(channel chan []byte,Request requestApi.Request) map[string]interface{} {
    result := make(map[string]interface{})
    return result
}

func WaitCloseLock(responseRabbit RequestRabbit, channel chan []byte,tid string , tidInt int,Request requestApi.Request) bool {
    keyAutomat := Request.Config.AutomatId + Request.Config.AccountId
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = "Замок открыт, заберите продукты"
    parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
    parametrs["status"] = transactionStruct.TransactionState_VendSession
    SetDataTransaction(parametrs,where)
    log.Println("2 WaitRabbit")
    responseRabbit = WaitResponseRabbit(channel,Bank,Request)
    if responseRabbit.St != IceboxStatus_End {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = TransactionStruct.GetDescriptionCodeCooler(responseRabbit.St)
        parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct)
        AddDataAutomatEvent(TransactionStruct)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    //log.Printf("\n [x] %v",transactionStruct.TransactionState_VendOk)
    // where["id"] = Request.IdTransaction
    // parametrs["ps_desc"] = "Замок закрыт. Спасибо за покупки"
    // parametrs["error"] = "Замок закрыт. Спасибо за покупки"
    // parametrs["status"] = transactionStruct.TransactionState_VendOk
    // SetDataTransaction(parametrs,where)
    // AddDataTransactionProducts(TransactionStruct)
    return true
}

func (cl Cooler) StartPrepayment(channel chan []byte,Request requestApi.Request) bool {
    keyAutomat := Request.Config.AutomatId + Request.Config.AccountId
    tidInt,_ := strconv.Atoi(Request.IdTransaction)
    tid := fmt.Sprintf("%v",Request.IdTransaction)
    resultMessage := make(map[string]interface{})
    resultMessage["tid"] = Request.IdTransaction
    resultMessage["sum"] = Request.Sum
    resultMessage["m"] = 4
    //resultMessage["products"] = TransactionStruct.Products
    resultMessage["a"] = 1
    data,err := json.Marshal(resultMessage)
    if err != nil {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = "ошибка преобразования map[string]interface{} в json (Отправка данных устройству)"
        parametrs["error"] = fmt.Sprintf("%s",err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = "Оплата успешна, ожидание открытия замка"
    parametrs["error"] = "Оплата успешна, ожидание открытия замка"
    parametrs["status"] = transactionStruct.TransactionState_VendSession
    SetDataTransaction(parametrs,where)
    SendMessageToDevice(data,Bank,Request)
    log.Printf("\n [x] %s","1 WaitRabbit")
    responseRabbit := WaitResponseRabbit(channel,Bank,Request)
    switch responseRabbit.St {
		 case IceboxStatus_Drink:
            return WaitCloseLock(responseRabbit,channel,tid,tidInt,Request)
		 case IceboxStatus_Icebox:
            return WaitCloseLock(responseRabbit,channel,tid,tidInt,Request)
         default:
            where["id"] = Request.IdTransaction
            parametrs["ps_desc"] = TransactionStruct.GetDescriptionCodeCooler(responseRabbit.St)
            parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
            parametrs["status"] = transactionStruct.TransactionState_Error
            SetDataTransaction(parametrs,where)
            AddDataTransactionProducts(TransactionStruct)
            AddDataAutomatEvent(TransactionStruct)
            UpdateAutomat(TransactionStruct)
            Bank.ReturnMoney(tid,Request)
            Transactions.RemoveChannel(tidInt)
            Transactions.RemoveReplayProtection(keyAutomat)
            return false
	 }
}

func (cl Cooler) InitDeviceData(data transactionStruct.Transaction){

}
/* Оплаты здесь нет. Вся оплата на уровне выше. Реализуем бизнес-логику общения с модемам в зависимости от двух типов: 
    1) PayType - тип оплаты 
    2) DeviceType - какой тип устройста, с которым сейчас идёт общение
*/
var connectDb connectionPostgresql.DatabaseInstance
var Rabbit *ConnectionRabbitMQ.ChannelMQ
var Config *configEphor.Config
var Bank interfaseBank.Bank
var Transactions transactionDispetcher.TransactionDispetcher
var TransactionStruct transactionStruct.Transaction
var TransactionProductStore transactionProduct.StoreTransactionProduct
var AutomatEventStore automatEvent.StoreAutomatEvent
var AutomatStore automat.StoreAutomat
var (
    where = make(map[string]interface{})
    parametrs = make(map[string]interface{})
    resultAnswerBank = make(map[string]interface{})
)

func (cl Cooler) SendMessage(data transactionStruct.Transaction,conn connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,conf *configEphor.Config, req requestApi.Request,bank interfaseBank.Bank,resultBank map[string]interface{},transactions transactionDispetcher.TransactionDispetcher,channel chan []byte) map[string]interface{} {
    TransactionProductStore.Connection = conn
    AutomatEventStore.Connection = conn
    AutomatStore.Connection = conn
    connectDb = conn
    Rabbit = rabbitmq
    Config = conf
    Request := req
    Bank = bank
    resultAnswerBank = resultBank
    Transactions = transactions
    TransactionStruct = data
    result := make(map[string]interface{})
    if data.PayType == transactionStruct.Prepayment {
        log.Printf("\n [x] %s","Prepayment")
        status := cl.StartPrepayment(channel,Request)
        result["status"] = status
        return result
    }
    if data.PayType == transactionStruct.Postpaid {
        status := cl.StartPostpaid(channel,Request)
        result["status"] = status
        return result
    }
    result["status"] = false
    return result

}

var trasactionData transactionStruct.Transaction

func (newCl *NewCoolerStruct) NewDevice() interfaceDevice.Device  /* тип interfaceDevice.Device*/ {
    return &NewCoolerStruct{
        Cooler: Cooler{
        Name: "Cooler",
        Status: 0,
       },
    }
}