package paymentmanager

import (
    factoryBank "ephorservices/src/server/utils/factory/payments"
    interfaseBank "ephorservices/src/server/utils/interface/payment"
    requestApi "ephorservices/src/data/requestApi"
    "fmt"
    "log"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    transactionStruct "ephorservices/src/data/transaction"
)

const (
    SBER = 1
    VENDPAY = 2
)

type PaymentManager struct{

}

func (pm PaymentManager) InitBank(typeBank int) interfaseBank.Bank {
    return factoryBank.GetBank(typeBank)
}

func (pm PaymentManager) StartPrepayment(request requestApi.Request, Bank interfaseBank.Bank) map[string]interface{} {
    result := make(map[string]interface{})
    Bank.InitBankData(connect)
    resultHold := Bank.HoldMoney(request)
    if resultHold["status"] == false {
        result["status"]    = transactionStruct.TransactionState_Error
        result["orderId"]   = fmt.Sprintf("%v",resultHold["orderId"])
        result["ps_invoice_id"] = fmt.Sprintf("%v", resultHold["invoiceId"])
        result["success"]   = false
        result["ps_desc"]   = resultHold["description"]
        result["error"]     = resultHold["message"]
        return result
    }
    orderId := fmt.Sprintf("%v",resultHold["orderId"])
    invoiceId := fmt.Sprintf("%v", resultHold["invoiceId"])
    log.Println(orderId)
    request.OrderId = orderId
    log.Printf("%+v",request)
    resultHold = Bank.DebitHoldMoney(orderId,0,request)
    if resultHold["status"] == false {
        result["status"]    = transactionStruct.TransactionState_Error
        result["orderId"]   = orderId
        result["ps_invoice_id"] = invoiceId
        result["success"]   = false
        result["ps_desc"]   = resultHold["description"]
        result["error"]     = resultHold["message"]
        return result
    }
    result["ps_desc"] = resultHold["description"]
    result["ps_order"] = orderId
    result["ps_invoice_id"] = invoiceId
    result["status"] = transactionStruct.TransactionState_MoneyDebit
    result["success"] = true
    return result
}

func (pm PaymentManager) StartPostpaid(request requestApi.Request,Bank interfaseBank.Bank) map[string]interface{} {
    result := make(map[string]interface{})
    Bank.InitBankData(connect)
    resultHold := Bank.HoldMoney(request)
    if resultHold["status"] == false {
        result["status"]    = transactionStruct.TransactionState_Error
        result["orderId"]   = nil
        result["success"]   = false
        result["ps_desc"]   = resultHold["description"]
        result["error"]     = resultHold["message"]
        return result
    }
    log.Printf("%+v",resultHold)
    result["ps_desc"] = resultHold["description"]
    result["ps_order"] = resultHold["orderId"]
    result["ps_invoice_id"] = resultHold["invoiceId"]
    result["status"] = transactionStruct.TransactionState_MoneyDebitOk
    result["success"] = true
    return result
}
var connect connectionPostgresql.DatabaseInstance

func (pm PaymentManager) StartСommunicationBank(request requestApi.Request,transaction transactionStruct.Transaction,connectBd connectionPostgresql.DatabaseInstance ) (map[string]interface{},interfaseBank.Bank) {
    result := make(map[string]interface{})
    payType := request.Config.PayType
    Bank := pm.InitBank(request.Config.BankType)
    connect = connectBd
    if Bank == nil {
        result["status"] = transactionStruct.TransactionState_Error
        result["success"] = false
        result["ps_desc"] = "no available bank"
        result["error"] = "no available bank"
        return result,nil
    }
    if payType == transactionStruct.Prepayment {
        result = pm.StartPrepayment(request,Bank)
        return result,Bank
    }
    if payType == transactionStruct.Postpaid {
        result = pm.StartPostpaid(request,Bank)
        return result,Bank
    }
    result["status"] = transactionStruct.TransactionState_Error
    result["success"] = false
    result["ps_desc"] = "no available paymentType. Awalable 1 or 2"
    result["error"] = "no available paymentType. Awalable 1 or 2"
    return result,Bank
}