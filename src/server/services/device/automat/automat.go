package automat 

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
    ware "ephorservices/src/server/model/schema/account/ware"
)


type Automat struct {
    Name string
    Status int
}

type NewAutomatStruct struct {
    Automat
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

func AddDataTransactionProducts(transaction transactionStruct.Transaction,Request requestApi.Request){
    params := make(map[string]interface{})
    if len(transaction.Products) < 1 {
        reqWare := make(map[string]interface{})
        reqWare["id"] = Request.WareId
        reqWare["account_id"] = transaction.AccountId
        ware := WareStore.GetWithOptions(reqWare)
        log.Printf("\n [x] %+v",ware)
        if len(ware) < 1 {
            params["transaction_id"] = transaction.Tid
            params["ware_id"] = Request.WareId
            params["value"] = Request.Sum
            params["quantity"] = 1
            TransactionProductStore.AddByParams(params)
            return
        }
        wareModel := ware[0]
        log.Printf("\n [x] %+v",wareModel)
        params["transaction_id"] = transaction.Tid
        params["ware_id"] = wareModel.Id
        params["value"] = Request.Sum
        params["name"],_ = wareModel.Name.Value()
        params["quantity"] = 1
        log.Printf("\n [x] %+v",params)
        TransactionProductStore.AddByParams(params)
        return
    }
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
    fmt.Println(len(transaction.Products))
    addSeconds := 1
    for _, product := range transaction.Products {
        date := GetDateTime(transaction.Date,addSeconds)
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

func SetDataTransaction(parametrs,where map[string]interface{})  {
    connectDb.Set("transaction",parametrs,where)
    connectDb.SetLog(fmt.Sprintf("%q",parametrs["error"]))
}

func SendMessageToDevice(jsonData []byte,bank interfaseBank.Bank, Request requestApi.Request) bool {
     tid := fmt.Sprintf("%v",trasactionData.Tid)
     err := Rabbit.PublishMessage(jsonData,fmt.Sprintf("ephor.1.dev.%v",Request.Imei))
     if err != nil {
         bank.ReturnMoney(tid,Request)
         return false
     }
     return true
}

func WaitResponseRabbit(channel chan []byte,Bank interfaseBank.Bank, Request requestApi.Request) RequestRabbit {
    timer := time.NewTimer(time.Duration(Config.RabbitMq.ExecuteTimeSeconds) * time.Second)
    tid := fmt.Sprintf("%v",trasactionData.Tid)
    log.Printf("\n [x] %s","Timer")
    requestRabbit := RequestRabbit{}
    select {
        case <-timer.C:{
            requestRabbit = RequestRabbit{}
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
            requestRabbit = RequestRabbit{}
            timer.Stop()
            log.Printf("\n [x] %s","Timer ok")
            log.Println("Timer ok")
            json.Unmarshal(result, &requestRabbit)
        }
    }
    return requestRabbit
}

func WaitStatusDevice(responseRabbit RequestRabbit,orderId,tid string,tidInt,keyAutomat int, Request requestApi.Request) bool {
    
    switch responseRabbit.St {
		 case transactionStruct.VendState_Approving,transactionStruct.VendState_Vending:
            Request.SumOneProduct = responseRabbit.Sum
            resultDebitMoney := Bank.DebitHoldMoney(orderId,responseRabbit.Sum,Request)
            if resultDebitMoney["status"] == false {
                where["id"] = Request.IdTransaction
                parametrs["ps_desc"] = "Не удалось списать денньги. Деньги вернутся в течении суток."
                parametrs["error"] = resultDebitMoney["message"]
                parametrs["status"] = transactionStruct.TransactionState_Error
                SetDataTransaction(parametrs,where)
                AddDataTransactionProducts(TransactionStruct,Request)
                Bank.ReturnMoney(tid,Request)
                Transactions.RemoveChannel(tidInt)
                Transactions.RemoveReplayProtection(keyAutomat)
                return false
            }
            StatusBank = 1
            where["id"] = Request.IdTransaction
            parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
            parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
            parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_Vending)
            SetDataTransaction(parametrs,where)
            return true
         case transactionStruct.VendState_VendOk:
            Request.SumOneProduct = responseRabbit.Sum
            if StatusBank == 1 {
                resultDebitMoney := Bank.DebitHoldMoney(orderId,responseRabbit.Sum,Request)
                if resultDebitMoney["status"] == false {
                    where["id"] = Request.IdTransaction
                    parametrs["ps_desc"] = "Не удалось списать денньги. Деньги вернутся в течении суток."
                    parametrs["error"] = resultDebitMoney["message"]
                    parametrs["status"] = transactionStruct.TransactionState_Error
                    SetDataTransaction(parametrs,where)
                    AddDataTransactionProducts(TransactionStruct,Request)
                    Bank.ReturnMoney(tid,Request)
                    Transactions.RemoveChannel(tidInt)
                    Transactions.RemoveReplayProtection(keyAutomat)
                    return false
                }
            }
            where["id"] = Request.IdTransaction
            parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
            parametrs["error"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
            parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_VendOk)
            SetDataTransaction(parametrs,where)
            AddDataTransactionProducts(TransactionStruct,Request)
            Transactions.RemoveChannel(tidInt)
            Transactions.RemoveReplayProtection(keyAutomat)
            return true
         default:
            where["id"] = Request.IdTransaction
            parametrs["ps_desc"] = "Деньги возвращены"
            parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
            parametrs["status"] = transactionStruct.TransactionState_Error
            SetDataTransaction(parametrs,where)
            AddDataTransactionProducts(TransactionStruct,Request)
            Bank.ReturnMoney(tid,Request)
            Transactions.RemoveChannel(tidInt)
            Transactions.RemoveReplayProtection(keyAutomat)
            return false
    }
}

 
func (au Automat) StartPostpaid(channel chan []byte,Request requestApi.Request) bool {
    keyAutomat := Request.Config.AutomatId + Request.Config.AccountId
    orderId := fmt.Sprintf("%v",resultAnswerBank["orderId"])
    tidInt,_ := strconv.Atoi(Request.IdTransaction)
    tid := fmt.Sprintf("%v",Request.IdTransaction)
    resultMessage := make(map[string]interface{})
    resultMessage["tid"] = Request.IdTransaction
    resultMessage["sum"] = Request.Sum
    resultMessage["m"] = 1
    resultMessage["a"] = 2
    data,err := json.Marshal(resultMessage)
    if err != nil {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = "ошибка преобразования map[string]interface{} в json (Отправка данных устройству)"
        parametrs["error"] = fmt.Sprintf("%s",err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    SendMessageToDevice(data,Bank,Request)
    responseRabbit := WaitResponseRabbit(channel,Bank,Request)
    if responseRabbit.St != transactionStruct.VendState_Session {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
        parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct,Request)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = "Оплата успешна, ожидание нажатия пользователем кнопки на ТА"
    parametrs["error"] = "Оплата успешна, ожидание нажатия пользователем кнопки на ТА"
    parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_Session)
    SetDataTransaction(parametrs,where)
    responseRabbit = WaitResponseRabbit(channel,Bank,Request)
    Request.WareId = strconv.Itoa(responseRabbit.Wid)
    Request.Sum = responseRabbit.Sum
    resultDevice := WaitStatusDevice(responseRabbit,orderId,tid,tidInt,keyAutomat,Request)
    if resultDevice == false {
        return false
    }
    responseRabbit = WaitResponseRabbit(channel,Bank,Request)
    resultDevice = WaitStatusDevice(responseRabbit,orderId,tid,tidInt,keyAutomat,Request)
    if resultDevice == false {
        return false
    }
    return true
}

func (au Automat) StartPrepayment(channel chan []byte,Request requestApi.Request) bool {
    keyAutomat := Request.Config.AutomatId + Request.Config.AccountId
    tidInt,_ := strconv.Atoi(Request.IdTransaction)
    tid := fmt.Sprintf("%v",Request.IdTransaction)
    resultMessage := make(map[string]interface{})
    resultMessage["tid"] = Request.IdTransaction
    resultMessage["sum"] = Request.Sum
    resultMessage["wid"] = Request.Products[0]["ware_id"]
    resultMessage["m"] = 1
    resultMessage["a"] = 1
    data,err := json.Marshal(resultMessage)
    if err != nil {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = "ошибка преобразования map[string]interface{} в json (Отправка данных устройству)"
        parametrs["error"] = fmt.Sprintf("%s",err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct,Request)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    SendMessageToDevice(data,Bank,Request)
    responseRabbit := WaitResponseRabbit(channel,Bank,Request)
    if responseRabbit.St != transactionStruct.VendState_Session {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
        parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct,Request)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = "Оплата успешна, ожидание нажатия пользователем кнопки на ТА"
    parametrs["error"] = "Оплата успешна, ожидание нажатия пользователем кнопки на ТА"
    parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_Session)
    SetDataTransaction(parametrs,where)
    responseRabbit = WaitResponseRabbit(channel,Bank,Request)
    if responseRabbit.St != transactionStruct.VendState_Vending {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
        parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct,Request)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
    parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
    parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_Vending)
    SetDataTransaction(parametrs,where)
    responseRabbit = WaitResponseRabbit(channel,Bank,Request)
    if responseRabbit.St != transactionStruct.VendState_VendOk {
        where["id"] = Request.IdTransaction
        parametrs["ps_desc"] = "Деньги возвращены"
        parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(TransactionStruct,Request)
        Bank.ReturnMoney(tid,Request)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return false
    }
    where["id"] = Request.IdTransaction
    parametrs["ps_desc"] = TransactionStruct.GetDescriptionCode(responseRabbit.St)
    parametrs["error"] = TransactionStruct.GetDescriptionErr(responseRabbit.Err)
    parametrs["status"] = TransactionStruct.GetStatusServer(transactionStruct.VendState_VendOk)
    SetDataTransaction(parametrs,where)
    AddDataTransactionProducts(TransactionStruct,Request)
    Transactions.RemoveChannel(tidInt)
    Transactions.RemoveReplayProtection(keyAutomat)
    return true
}

func (au Automat) InitDeviceData(data transactionStruct.Transaction){

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
var WareStore ware.StoreWare
var (
    where = make(map[string]interface{})
    parametrs = make(map[string]interface{})
    resultAnswerBank = make(map[string]interface{})
    StatusBank = 0 
)

func (au Automat) SendMessage(data transactionStruct.Transaction,conn connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,conf *configEphor.Config, req requestApi.Request,bank interfaseBank.Bank,resultBank map[string]interface{},transactions transactionDispetcher.TransactionDispetcher,channel chan []byte) map[string]interface{} {
    TransactionProductStore.Connection = conn
    AutomatEventStore.Connection = conn
    connectDb = conn
    WareStore.Connection = conn
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
        status := au.StartPrepayment(channel,Request)
        result["status"] = status
        return result
    }
    if data.PayType == transactionStruct.Postpaid {
        status := au.StartPostpaid(channel,Request)
        result["status"] = status
        return result
    }
    result["status"] = false
    return result

}

var trasactionData transactionStruct.Transaction

func (newAu NewAutomatStruct) NewDevice() interfaceDevice.Device  /* тип interfaceDevice.Device*/ {
    return &NewAutomatStruct{
        Automat: Automat{
        Name: "Automat",
        Status: 0,
       },
    }
}