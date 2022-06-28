package transaction

import(
    "fmt"
    "time"
    "log"
    "runtime"
    "strconv"
    "math/rand"
    config "ephorservices/src/configs"
    transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
    deviceFactory "ephorservices/src/server/utils/factory/device"
    deviceInterfase "ephorservices/src/server/utils/interface/device"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    transactionStruct "ephorservices/src/data/transaction"
    paymentmanager "ephorservices/src/server/utils/paymentmanager"
    transactionProduct "ephorservices/src/server/model/schema/main/transactionproduct"
    automatEvent "ephorservices/src/server/model/schema/account/automatevent" 
    transactionStore "ephorservices/src/server/model/schema/main/transaction"
    requestApi "ephorservices/src/data/requestApi"
    fiscalmanager "ephorservices/src/server/middleLayer/fiscal"
    automatlocation "ephorservices/src/server/model/schema/account/automatlocation"
    automat "ephorservices/src/server/model/schema/account/automat"
)

func Finish(){
    runtime.Goexit()
}

func GetDateTime(stringTime string,seconds int) string {
    t,_ := time.Parse("2006-01-02 15:04:05",stringTime)
    t2 := t.Add(-180 * time.Minute)
    newTime := t2.Add(time.Duration(seconds) * time.Second)
    resultTime := newTime.Format("2006-01-02 15:04:05")
    return resultTime
}

func GetAutomatLocation(idAutomat,accountId interface{}) ([]map[string]interface{}) {
   options := make(map[string]interface{})
   options["account_id"] = accountId
   options["automat_id"] = idAutomat
   options["to_date"] = nil
   return AutomatLocationStore.GetWithOptions(options)
}

func GetAutomat(id,accountId interface{}) ([]map[string]interface{},bool) {
   return AutomatStore.GetOneById(id,accountId);
}

func UpdateAutomat(transaction transactionStruct.Transaction){
    log.Println("UpdateAutomat")
    automat,_ := GetAutomat(transaction.AutomatId,transaction.AccountId)
    if len(automat) < 1 {
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

func InitDevice(device int) deviceInterfase.Device {
    return deviceFactory.GetDevice(device)
}

func SetDataTransactionProducts(transactionProduct map[string]interface{}){
    TransactionProductStore.SetByParams(transactionProduct)
}

func AddDataTransactionProducts(transaction transactionStruct.Transaction){
    params := make(map[string]interface{})
    for _, product := range transaction.Products {
        params["transaction_id"] = transaction.Tid
        params["name"] = product["name"]
        params["select_id"] = product["select_id"]
        params["ware_id"] = product["ware_id"]
        params["value"] = product["price"]
        params["tax_rate"] = product["tax_rate"]
        params["quantity"] = product["quantity"]
        TransactionProductStore.AddByParams(params) 
    }
}

func Random(min int, max int) int {
    return rand.Intn(max-min) + min
}

func AddDataAutomatEvent(transaction transactionStruct.Transaction,frResult map[string]interface{}){
    params := make(map[string]interface{})
    addSeconds := 1
    for _, product := range transaction.Products {
        for i := 0; i < int(product["quantity"].(float64)); i++ {
            date := GetDateTime(transaction.Date,addSeconds)
            params["account_id"] = transaction.AccountId
            params["automat_id"] = transaction.AutomatId
            params["type"] = 3
            params["date"] = date
            params["category"] = 1
            params["name"] = product["name"]
            params["credit"] = product["price"]
            params["select_id"] = product["select_id"]
            params["ware_id"] = product["ware_id"]
            params["value"] = product["price"]
            params["status"] = frResult["fr_status"]
            params["error_detail"] = "нет ошибок"
            params["type_fr"] = frResult["type_fr"]
            params["fp"] = frResult["fp"]
            params["fn"] = frResult["fn"]
            params["fd"] = frResult["fd"]
            params["payment_device"] = product["payment_device"]
            params["tax_system"] = transaction.Tax_system
            params["tax_rate"] = product["tax_rate"]
            params["tax_value"] = product["tax_value"]
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

func SetDataTransaction(parametrs,where map[string]interface{})  {
    connectDb.Set("transaction",parametrs,where)
    connectDb.SetLog(fmt.Sprintf("%q",parametrs["error"]))
}

func initDataTransaction(request requestApi.Request) transactionStruct.Transaction {
    products := request.Products
    transaction := transactionStruct.Transaction{}
    transaction.Tid = request.IdTransaction
    transaction.Date = request.Date
    transaction.Sum = request.Sum
    transaction.Token = request.PaymentToken
    transaction.PayType = request.Config.PayType
    transaction.AccountId = request.Config.AccountId
    transaction.AutomatId = request.Config.AutomatId
    transaction.DeviceType = request.Config.DeviceType
    transaction.UserPhone = request.Config.UserPhone
	transaction.ReturnUrl = request.Config.ReturnUrl
	transaction.DeepLink = request.Config.DeepLink
	transaction.TokenType = request.Config.TokenType
    transaction.SumMax = request.SumMax
    transaction.QrFormat = request.Config.QrFormat
    for _, product := range products {
        transaction.Products = append(transaction.Products,product)
    }
    return transaction
}

func CheckTransactionOfAutomat(request requestApi.Request) bool {
    keyFound := request.Config.AutomatId + request.Config.AccountId
    resultFound := Transactions.GetReplayProtection(keyFound)
    if resultFound == false {
        return false
    }
    return true
}

func StartTransaction(bankChannel chan bool, json_data requestApi.Request){
    where = make(map[string]interface{})
    parametrs = make(map[string]interface{})
    request := json_data
    tidInt,_ := strconv.Atoi(request.IdTransaction)
    keyAutomat := request.Config.AutomatId + request.Config.AccountId
    transactionData := initDataTransaction(request)
    automatLocationModel :=  GetAutomatLocation(transactionData.AutomatId,transactionData.AccountId)
    if len(automatLocationModel) != 0 {
       transactionData.PointId = int(automatLocationModel[0]["company_point_id"].(int64))
    }
    result,Bank := PaymentManager.StartСommunicationBank(request,transactionData,connectDb)
    if result["success"] == false {
        where["id"] = request.IdTransaction
        parametrs["ps_order"] = result["ps_order"]
        parametrs["ps_invoice_id"] = result["ps_invoice_id"]
        parametrs["ps_desc"] = result["ps_desc"]
        parametrs["error"] = result["error"]
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(transactionData)
        return 
    }
    request.OrderId = string(result["ps_order"].(string))
    where["id"] = request.IdTransaction
    parametrs["ps_order"] = result["ps_order"]
    parametrs["ps_invoice_id"] = result["ps_invoice_id"]
    parametrs["ps_desc"] = result["ps_desc"]
    parametrs["error"] = result["error"]
    parametrs["status"] = transactionStruct.TransactionState_MoneyDebitOk
    SetDataTransaction(parametrs,where)
    channel := Transactions.AddChannel(tidInt)
    Transactions.AddReplayProtection(keyAutomat,request.Config.AutomatId)
    device := InitDevice(transactionData.DeviceType)
    if device == nil {
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        where["id"] = request.IdTransaction 
        parametrs["error"] = "no available device type"
        parametrs["status"] = transactionStruct.TransactionState_Error
        SetDataTransaction(parametrs,where)
        AddDataTransactionProducts(transactionData)
        Bank.ReturnMoney(request.IdTransaction,request)
        return 
    }
    resultSendMessageDevice := device.SendMessage(transactionData,connectDb,Rabbit,conf,request,Bank,result,*Transactions,channel)
    if resultSendMessageDevice["status"] == false {
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        return
    }
    if transactionData.DeviceType == 7 {
        UpdateAutomat(transactionData)
        where["id"] = request.IdTransaction 
        parametrs["error"] = "Пробиваем чек, подождите"
        parametrs["status"] = transactionStruct.TransactionState_WaitFiscal
        parametrs["ps_desc"] = "Пробиваем чек, подождите"
        SetDataTransaction(parametrs,where)
        resultFr,transactionFiscal := FiscalManager.FiscalProducts(transactionData)
        if resultFr["status"] == false {
            Transactions.RemoveChannel(tidInt)
            Transactions.RemoveReplayProtection(keyAutomat)
            where["id"] = request.IdTransaction 
            parametrs["error"] = resultFr["message"]
            parametrs["f_status"] = resultFr["fr_status"]
            parametrs["status"] = transactionStruct.TransactionState_VendOk
            parametrs["ps_desc"] = "Спасибо за покупку"
            SetDataTransaction(parametrs,where)
            AddDataTransactionProducts(transactionFiscal)
            return
        }
        where["id"] = request.IdTransaction 
        parametrs["error"] = "Нет ошибок"
        parametrs["status"] = transactionStruct.TransactionState_VendOk
        parametrs["ps_desc"] = "Спасибо за покупку"
        SetDataTransaction(parametrs,where)
        Transactions.RemoveChannel(tidInt)
        Transactions.RemoveReplayProtection(keyAutomat)
        AddDataTransactionProducts(transactionData)
        AddDataAutomatEvent(transactionFiscal,resultFr)
    }
    return 
}

func GateDateTimeMoscow() string {
    t := time.Now()
    t.Add(-72 * time.Hour)
    resultTime := t.Format("2006-01-02 15:04:05")
    return resultTime
}

func AddTransaction(request requestApi.Request) (string,int) {
    rand.Seed(time.Now().UnixNano())
    randNoise := Random(10000000,20000000)
    parametrInsert := make(map[string]interface{})
    parametrInsert["automat_id"] = request.Config.AutomatId
    parametrInsert["account_id"] = request.Config.AccountId
    parametrInsert["token_id"] = request.PaymentToken
    parametrInsert["status"] = 1
    parametrInsert["date"] = GateDateTimeMoscow()
    parametrInsert["pay_type"] = request.Config.PayType
    parametrInsert["sum"] = request.Sum
    parametrInsert["ps_type"] = request.Config.BankType
    parametrInsert["token_type"] = request.Config.TokenType
    parametrInsert["qr_format"] = request.Config.QrFormat
    Id := TransactionStore.AddByParams(parametrInsert)
    noise := randNoise + Id
    updateNoise := make(map[string]interface{})
    updateNoise["id"] = Id
    updateNoise["noise"] = strconv.Itoa(noise)
    TransactionStore.SetByParams(updateNoise)
    return strconv.Itoa(noise),Id
}


var (
    connectDb connectionPostgresql.DatabaseInstance
    conf *config.Config
    Rabbit *ConnectionRabbitMQ.ChannelMQ
    Transactions *transactionDispetcher.TransactionDispetcher
    FiscalManager fiscalmanager.FiscalMiddleLayer
    TransactionStruct transactionStruct.Transaction
    PaymentManager paymentmanager.PaymentManager
    TransactionProductStore transactionProduct.StoreTransactionProduct
    AutomatEventStore automatEvent.StoreAutomatEvent
    TransactionStore  transactionStore.StoreTransaction
    AutomatLocationStore automatlocation.StoreAutomatLocation
    AutomatStore automat.StoreAutomat
    where  map[string]interface{}
    parametrs map[string]interface{}
)

func Init(cfg *config.Config,connectPg connectionPostgresql.DatabaseInstance,rabbitmq *ConnectionRabbitMQ.ChannelMQ,transactions *transactionDispetcher.TransactionDispetcher){
    TransactionProductStore.Connection = connectPg
    AutomatEventStore.Connection = connectPg
    TransactionStore.Connection = connectPg
    AutomatLocationStore.Connection = connectPg
    AutomatStore.Connection = connectPg
    connectDb = connectPg
    conf = cfg
    Rabbit = rabbitmq
    Transactions = transactions
}



