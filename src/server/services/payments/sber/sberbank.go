package sber

import(
    "bytes"
	"encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "math"
    "time"
    randString "ephorservices/src/pkg/randgeneratestring"
    interfaceBank "ephorservices/src/server/utils/interface/payment"
    requestApi "ephorservices/src/data/requestApi"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    transactionStruct "ephorservices/src/data/transaction"
)
const TransactionState_Idle 				= 0; // Transaction Idle
const TransactionState_MoneyHoldStart 		= 1; // создали транзакцию банка
const TransactionState_MoneyHoldWait 		= 2; // ожидает ответ от банка
const TransactionState_MoneyDebitStart		= 8;
const TransactionState_MoneyDebitWait		= 9;
const TransactionState_MoneyDebitOk			= 10;
const TransactionState_Error 				= 120;
const TransactionState_SberPay				= 13;

const Order_Registered = 0 //- заказ зарегистрирован, но не оплачен;
const Order_HoldMoney  = 1 //- предавторизованная сумма удержана (для двухстадийных платежей);
const Order_FullAuthorizationOfTheAmount =  2 //- проведена полная авторизация суммы заказа;
const Order_AuthorizationCanceled = 3 //- авторизация отменена;
const Order_RefundOperationPerformed  = 4 //- по транзакции была проведена операция возврата;
const Order_AuthorizationThroughTheServerHasBeenInitiated = 5 //- инициирована авторизация через сервер контроля доступа банка-эмитента;
const Order_AuthorizationDenied = 6 //- авторизация отклонена.

type Sber struct {
    Name string
    Counter int
    PaymentType int
    UrlCreateOrder string
    UrlGetStatusOrder string
    UrlDepositeOrder string
    UrlReverse string
    UrlCreateOrderSberPay string
    UrlGetStatusOrderSberPay string
    UrlDepositeOrderSberPay string
    UrlReverseSberPay string
    UrlRefundSberPay string
    ConnectBD connectionPostgresql.DatabaseInstance
    Status int
    Req requestApi.Request
    Res interfaceBank.Response
}

type NewSberStruct struct {
    Sber
}

var checkTypeCheck = 0

func SetDataTransaction(parametrs,where map[string]interface{})  {
    ConnectBD.Set("transaction",parametrs,where)
}

func (s Sber) startTimeOut(checkType int,orderId string,sum int,Request requestApi.Request) (map[string]interface{}) {
    return s.CheckStatus(Request)
}

func (s Sber) MakeJsonRequestDepositOrderSberPay(sum int,Request requestApi.Request,orderId string) (string,error){
    stringRequest := ""
    stringRequest += fmt.Sprintf("userName=%s&",Request.Config.Login)
    stringRequest += fmt.Sprintf("password=%s&",Request.Config.Password)
    stringRequest += fmt.Sprintf("orderId=%s&",Request.OrderId)
    stringRequest += fmt.Sprintf("amount=%v",sum)
    return stringRequest ,nil
}

func (s Sber) MakeJsonRequestDepositOrder(sum int,Request requestApi.Request) ([]byte,error){
    requestOrder := make(map[string]interface{})
    requestOrder["amount"] = sum
    requestOrder["orderId"] = Request.OrderId
    data, err := json.Marshal(requestOrder)
    if err != nil {
        log.Printf("%+v",err)
        return nil, err
    }else {
        return data ,nil
    }
}

func (s Sber) MakeJsonRequestStatusOrder(Request requestApi.Request) ([]byte,error){
    requestOrder := make(map[string]interface{})
    requestOrder["token"] = Request.PaymentToken
    requestOrder["orderId"] = Request.OrderId
    data, err := json.Marshal(requestOrder)
    if err != nil {
        log.Printf("%+v",err)
        return nil, err
    }else {
        return data ,nil
    }
}

func (s Sber) MakeJsonRequestStatusOrderSberPay(Request requestApi.Request) (string,error){
    stringRequest := ""
    stringRequest += fmt.Sprintf("userName=%s&",Request.Config.Login)
    stringRequest += fmt.Sprintf("password=%s&",Request.Config.Password)
    stringRequest += fmt.Sprintf("orderId=%s",Request.OrderId)
    return stringRequest ,nil
}

func (s Sber) MakeJsonSberPayCreateOrder(Request requestApi.Request) (string, error){
    stringRequest := ""
    var orderString randString.GenerateString
    orderString.RandStringRunes()
    orderNumber := orderString.String
    orderNumber += fmt.Sprintf("%v",Request.Config.AutomatId)
    if Request.Config.ReturnUrl == "" {
        Request.Config.ReturnUrl = "https://paytest.ephor.online"
    }
    stringRequest += fmt.Sprintf("userName=%s&",Request.Config.Login)
    stringRequest += fmt.Sprintf("password=%s&",Request.Config.Password)
    stringRequest += fmt.Sprintf("orderNumber=%s&",orderNumber)
    stringRequest += fmt.Sprintf("amount=%v&",Request.Sum)
    stringRequest += fmt.Sprintf("returnUrl=%s&",Request.Config.ReturnUrl)
    stringRequest += fmt.Sprintf("description=%s&",Request.Config.Description)
    jsonParams := make(map[string]interface{}) 
    if Request.Config.TokenType == transactionStruct.TypeTokenSberPayAndroid || Request.Config.TokenType == transactionStruct.TypeTokenSberPayiOS {
       jsonParams["app2app"] = true
       jsonParams["app.deepLink"] = Request.Config.DeepLink
       if Request.Config.TokenType == transactionStruct.TypeTokenSberPayAndroid {
           jsonParams["app.osType"] = "android"
       }
       if Request.Config.TokenType == transactionStruct.TypeTokenSberPayiOS {
           jsonParams["app.osType"] = "ios"
       }
    }else if Request.Config.TokenType == transactionStruct.TypeTokenSberPayWeb {
        jsonParams["back2app"] = true
        stringRequest += fmt.Sprintf("phone=%s&",Request.Config.UserPhone)
    }
    
    data,_ := json.Marshal(jsonParams)
    stringRequest += fmt.Sprintf("jsonParams=%s",data)
    log.Printf("%s",stringRequest)
    return stringRequest ,nil
}

func (s Sber) MakeJsonOrderRequestCreateOrder(Request requestApi.Request) ([]byte, error){
    var orderString randString.GenerateString
    orderString.RandStringRunes()
    orderNumber := orderString.String
    requestOrder := make(map[string]interface{})
    requestOrder["merchant"] = Request.MerchantId
    requestOrder["orderNumber"] = orderNumber
    requestOrder["language"] = Request.Config.Language
    requestOrder["preAuth"] = true
    requestOrder["description"] = Request.Config.Description
    requestOrder["paymentToken"] = Request.PaymentToken
    requestOrder["amount"] = Request.Sum
    requestOrder["currencyCode"] = Request.Config.CurrensyCode
    requestOrder["returnUrl"]  = "https://test.ru" 
    log.Printf("%+v",requestOrder)
    data, err := json.Marshal(requestOrder)
    if err != nil {
        log.Printf("%+v",err)
        return nil, err
    }else {
        return data ,nil
    }
}

func (s Sber) CheckStatus(Request requestApi.Request) (map[string]interface{}) {
    count := 0
        for {
            select {
                case <-time.After(1 * time.Minute):
                    resultResponse,_ :=  s.getStatusSberPay(false,Request)
                    if resultResponse["status"] != true {
                        return resultResponse
                    }
                    ResponsePaymentSystem["status"] = false
                    ResponsePaymentSystem["message"] = "Нет ответа от банка"
                    ResponsePaymentSystem["description"] = "Нет ответа от банка"
                    ResponsePaymentSystem["orderId"] = fmt.Sprintf("%v",Request.OrderId)
                    ResponsePaymentSystem["invoiceId"] = fmt.Sprintf("%v",Request.SbolBankInvoiceId)
                    ResponsePaymentSystem["code"] = TransactionState_Error
                    return ResponsePaymentSystem
                case <-time.After(5 * time.Second):
                    resultResponse,Response :=  s.getStatusSberPay(false,Request)
                    count += 1
                    if resultResponse["status"] != true {
                        return resultResponse
                    }
                    if count >= 36 {
                        ResponsePaymentSystem["status"] = false
                        ResponsePaymentSystem["message"] = "Нет ответа от банка"
                        ResponsePaymentSystem["description"] = "Нет ответа от банка"
                        ResponsePaymentSystem["orderId"] = fmt.Sprintf("%v",Request.OrderId)
                        ResponsePaymentSystem["invoiceId"] = fmt.Sprintf("%v",Request.SbolBankInvoiceId)
                        ResponsePaymentSystem["code"] = TransactionState_Error
                        return ResponsePaymentSystem
                    }
                    if Response.ErrorCode !=0 {
                        ResponsePaymentSystem["status"] = false
                        ResponsePaymentSystem["message"] = Response.ErrorMessage
                        ResponsePaymentSystem["description"] = Response.ErrorMessage
                        ResponsePaymentSystem["code"] = TransactionState_Error
                        return ResponsePaymentSystem
                    }
                    if Response.OrderStatus != 1 {
                        if Response.OrderStatus == 3 {
                            ResponsePaymentSystem["status"] = false
                            ResponsePaymentSystem["message"] = "Авторизация отменена"
                            ResponsePaymentSystem["description"] = "Авторизация отменена"
                            ResponsePaymentSystem["code"] = TransactionState_Error
                            return ResponsePaymentSystem
                        }else if Response.OrderStatus == 4 {
                            ResponsePaymentSystem["status"] = false
                            ResponsePaymentSystem["message"] = "По транзакции была проведена операция возврата"
                            ResponsePaymentSystem["description"] = "По транзакции была проведена операция возврата"
                            ResponsePaymentSystem["code"] = TransactionState_Error
                            return ResponsePaymentSystem
                        } else if Response.OrderStatus == 6 {
                            ResponsePaymentSystem["status"] = false
                            ResponsePaymentSystem["message"] = "Авторизация отклонена"
                            ResponsePaymentSystem["description"] = "Авторизация отклонена"
                            ResponsePaymentSystem["code"] = TransactionState_Error
                            return ResponsePaymentSystem
                        }
                    }else {
                        where := make(map[string]interface{})
                        parametrs := make(map[string]interface{})
                        where["id"] = Request.IdTransaction
                        parametrs["ps_desc"] = "Платёж подтверждён"
                        parametrs["error"] = "Платёж подтверждён"
                        parametrs["status"] = Order_FullAuthorizationOfTheAmount
                        SetDataTransaction(parametrs,where)
                        log.Printf("%+v",Response)
                        ResponsePaymentSystem["status"] = true
                        ResponsePaymentSystem["message"] = "Платёж подтверждён"
                        ResponsePaymentSystem["description"] = "Платёж подтверждён"
                        ResponsePaymentSystem["orderId"] = fmt.Sprintf("%v",Request.OrderId)
                        ResponsePaymentSystem["invoiceId"] = fmt.Sprintf("%v",Request.SbolBankInvoiceId)
                        ResponsePaymentSystem["code"] = Order_FullAuthorizationOfTheAmount
                        return ResponsePaymentSystem
                    }
		    }
        }
}

func (s Sber) debitMoneySberPay(orderId string,sum int,Request requestApi.Request) (map[string]interface{},interfaceBank.Response) {
    Response := interfaceBank.Response{}
    if sum != 0 {Request.Sum = sum}
    dataPush,err := s.MakeJsonRequestDepositOrderSberPay(Request.Sum,Request,orderId)
    if err != nil {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
        ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Списание денег транзакции)"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem,Response
    }
    url := fmt.Sprintf("%s?%s",s.UrlDepositeOrderSberPay,dataPush)
    fmt.Printf("%s",url)
    Response = s.Call("POST",url,nil)
    if Response.ErrorCode != 0 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem,Response
    }else if Response.ErrorCode == 7 {
       return s.debitMoneySberPay(orderId,sum,Request)
    }
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Деньги списаны"
    ResponsePaymentSystem["description"] = "Деньги списаны"
    ResponsePaymentSystem["code"] = TransactionState_MoneyDebitOk
    return ResponsePaymentSystem,Response
}

func (s Sber) getStatusSberPay(flag bool,Request requestApi.Request) (map[string]interface{},interfaceBank.Response) {
    Response := interfaceBank.Response{}
    dataPush,err := s.MakeJsonRequestStatusOrderSberPay(Request)
    if err != nil {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
        ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Опрос статуса sberpay транзакции)"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem,Response
    }
    url := fmt.Sprintf("%s?%s",s.UrlGetStatusOrderSberPay,dataPush)
    Response = s.Call("POST",url,nil)
    if Response.ErrorCode !=0 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem,Response
    }
    fmt.Sprintf("%s",Response.OrderStatus)
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Ожидание оплаты"
    ResponsePaymentSystem["description"] = "Ожидание оплаты"
    ResponsePaymentSystem["code"] = TransactionState_SberPay
    return ResponsePaymentSystem,Response
}

func (s Sber) registerPreAuth(Request requestApi.Request) map[string]interface{} {
    dataPush,err := s.MakeJsonSberPayCreateOrder(Request)
    fmt.Printf("%+v",dataPush)
    if err != nil {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
        ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Создание транзакции)"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    url := fmt.Sprintf("%s?%s",s.UrlCreateOrderSberPay,dataPush)
    fmt.Printf("%s",url)
    Response := s.Call("POST",url,nil)
    fmt.Printf("\n%+v\n",Response)
    if Response.Code != 0 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    if len(Response.ExternalParams) < 1 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }  
    sbolBankInvoiceId := Response.ExternalParams["sbolBankInvoiceId"]
    sbolInactive := Response.ExternalParams["sbolInactive"]
    if sbolInactive == false {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = "Банк недоступен"
        ResponsePaymentSystem["description"] = "Банк недоступен"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    ResponsePaymentSystem["status"] = false
    ResponsePaymentSystem["message"] = "Заказ принят, ожидание оплаты"
    ResponsePaymentSystem["description"] = "Заказ принят, ожидание оплаты"
    ResponsePaymentSystem["code"] = TransactionState_SberPay
    ResponsePaymentSystem["orderId"] = fmt.Sprintf("%v",Response.OrderId)
    ResponsePaymentSystem["invoiceId"] = fmt.Sprintf("%v",sbolBankInvoiceId)
    where := make(map[string]interface{})
    parametrs := make(map[string]interface{})
    where["id"] = Request.IdTransaction
    parametrs["ps_order"] = ResponsePaymentSystem["orderId"]
    parametrs["ps_invoice_id"] = ResponsePaymentSystem["invoiceId"]
    parametrs["ps_desc"] = "Заказ принят, ожидание оплаты"
    parametrs["error"] = "Заказ принят, ожидание оплаты"
    parametrs["status"] = TransactionState_SberPay
    SetDataTransaction(parametrs,where)
    orderId := string(ResponsePaymentSystem["orderId"].(string))
    Request.OrderId = fmt.Sprintf("%v",Response.OrderId)
    Request.SbolBankInvoiceId = fmt.Sprintf("%v",sbolBankInvoiceId)
    log.Printf(Request.OrderId)
    log.Printf(Request.SbolBankInvoiceId)
    return s.startTimeOut(1,orderId,0,Request)
}
var ConnectBD connectionPostgresql.DatabaseInstance
func (s Sber) InitBankData(connect connectionPostgresql.DatabaseInstance){
    ConnectBD = connect
}

func (s *Sber) HoldMoney(Request requestApi.Request) map[string]interface{} {
    fmt.Printf("%+v",Request)
    fmt.Printf("%s",transactionStruct.TypeTokenSberPayAndroid)
    if Request.Config.TokenType == 4 || Request.Config.TokenType == 5 || Request.Config.TokenType == 6 {
       return s.registerPreAuth(Request)
    }else {
        dataPush,err := s.MakeJsonOrderRequestCreateOrder(Request)
        log.Printf("%+v",dataPush)
        if err != nil {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
            ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Создание транзакции)"
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }

        Response := s.Call("POST",s.UrlCreateOrder,dataPush)
        log.Printf("%+v",Response)
        if Response.Success == false {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = Response.Error.Description
            ResponsePaymentSystem["description"] = Response.Error.Description
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }
        ResponsePaymentSystem["orderId"] = Response.Data["orderId"]
        return s.getStatus(Request)
    }
}

func (s Sber) getStatus(Request requestApi.Request) map[string]interface{} {
    dataPush,err := s.MakeJsonRequestStatusOrder(Request)
    if err != nil {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
        ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Опрос статуса транзакции)"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }

   Response := s.Call("POST",s.UrlGetStatusOrder,dataPush)
    if Response.ActionCode != 0 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    sum := float64(Request.Sum)
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции",math.Round((sum/100)))
    ResponsePaymentSystem["description"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции",math.Round((sum/100)))
    ResponsePaymentSystem["code"] = Order_FullAuthorizationOfTheAmount
    return ResponsePaymentSystem
}

func (s Sber) DebitSberPay(orderId string,sum int,Request requestApi.Request) map[string]interface{} {
    resultResponse,_ := s.debitMoneySberPay(orderId,sum,Request)
    if resultResponse["status"] != true {
        return resultResponse
    }
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Деньги списаны"
    ResponsePaymentSystem["description"] = "Деньги списаны"
    ResponsePaymentSystem["code"] = TransactionState_MoneyDebitOk
    return ResponsePaymentSystem
}

func (s Sber) DebitHoldMoney(orderId string,sum int,Request requestApi.Request) map[string]interface{} {
    if sum != 0 {Request.Sum = sum}
    if Request.Config.TokenType == 4 || Request.Config.TokenType == 5 || Request.Config.TokenType == 6 {
         return s.DebitSberPay(orderId,sum,Request)
    } else {
         dataPush,err := s.MakeJsonRequestDepositOrder(Request.Sum,Request)
        if err != nil {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
            ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Списание денег транзакции)"
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }
        
        Response := s.Call("POST",s.UrlDepositeOrder,dataPush)
        if Response.ErrorCode !=0 {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = Response.ErrorMessage
            ResponsePaymentSystem["description"] = Response.ErrorMessage
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }
        ResponsePaymentSystem["status"] = true
        ResponsePaymentSystem["message"] = "Деньги списаны"
        ResponsePaymentSystem["description"] = "Деньги списаны"
        ResponsePaymentSystem["code"] = TransactionState_MoneyDebitOk
        s.Status = 1
        return ResponsePaymentSystem
    }
}
 
func (s Sber) ReturnMoneySber(Request requestApi.Request) map[string]interface{} {
    stringRequest := ""
    stringRequest += fmt.Sprintf("userName=%s&",Request.Config.Login)
    stringRequest += fmt.Sprintf("password=%s&",Request.Config.Password)
    stringRequest += fmt.Sprintf("orderId=%s&",Request.OrderId)
    stringRequest += fmt.Sprintf("amount=%v",Request.Sum)
    url := fmt.Sprintf("%s?%s",s.UrlReverseSberPay,stringRequest)
    Response := s.Call("POST",url,nil)
    if Response.ErrorCode != 0 {
        ResponsePaymentSystem["status"] = false
        ResponsePaymentSystem["message"] = Response.ErrorMessage
        ResponsePaymentSystem["description"] = Response.ErrorMessage
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Деньги возвращены"
    ResponsePaymentSystem["description"] = "Деньги возвращены"
    ResponsePaymentSystem["code"] = TransactionState_Error
    return ResponsePaymentSystem
}

func (s Sber) ReturnMoney(orderId string,Request requestApi.Request) map[string]interface{} {
     if Request.Config.TokenType == transactionStruct.TypeTokenSberPayAndroid || Request.Config.TokenType == transactionStruct.TypeTokenSberPayiOS || Request.Config.TokenType == transactionStruct.TypeTokenSberPayWeb {
        return s.ReturnMoneySber(Request)
    }else {
        requestOrder := make(map[string]interface{})
        requestOrder["orderId"] = orderId
        data, err := json.Marshal(requestOrder)
        if err != nil {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = fmt.Sprintf("%s",err)
            ResponsePaymentSystem["description"] = "ошибка преобразования map[string]interface{} в json (Возврат денег)"
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }
        Response := s.Call("POST",s.UrlReverse,data)
        if Response.ErrorCode != 0 {
            ResponsePaymentSystem["status"] = false
            ResponsePaymentSystem["message"] = Response.ErrorMessage
            ResponsePaymentSystem["description"] = Response.ErrorMessage
            ResponsePaymentSystem["code"] = TransactionState_Error
            return ResponsePaymentSystem
        }
        ResponsePaymentSystem["status"] = true
        ResponsePaymentSystem["message"] = "Деньги возвращены"
        ResponsePaymentSystem["description"] = "Деньги возвращены"
        ResponsePaymentSystem["code"] = TransactionState_Error
        return ResponsePaymentSystem
    }
    
}


func (s Sber) Timeout(){
}

func (s Sber) GetPaymentType(){

}

func (s Sber) Call(method string, url string, json_request []byte) interfaceBank.Response {
    Response := interfaceBank.Response{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		Response.Code = 0
		Response.Status = TransactionState_Error
	}
	body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("%s",body)
    log.Printf("%s",url)
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), &Response)
    log.Printf("%+v",Response)
    return Response
}

var connectDb connectionPostgresql.DatabaseInstance
var RabbitMQ ConnectionRabbitMQ.ChannelMQ
var ResponsePaymentSystem map[string]interface{}
func (sber NewSberStruct) NewBank() interfaceBank.Bank  /* тип interfaceBank.Bank*/ {
    ResponsePaymentSystem = make(map[string]interface{})
    return &NewSberStruct{
        Sber: Sber{
        Name: "Sber",
        Counter: 0,
        PaymentType: 1, // srandart type payment 
        Status: 0,
        UrlCreateOrder: "https://3dsec.sberbank.ru/payment/google/payment.do",
        UrlGetStatusOrder: "https://3dsec.sberbank.ru/payment/google/getOrderStatusExtended.do",
        UrlDepositeOrder: "https://3dsec.sberbank.ru/payment/google/deposit.do", 
        UrlReverse: "https://3dsec.sberbank.ru/payment/google/reverse.do",
        UrlCreateOrderSberPay: "https://securepayments.sberbank.ru/payment/rest/registerPreAuth.do",
        UrlGetStatusOrderSberPay: "https://securepayments.sberbank.ru/payment/rest/getOrderStatusExtended.do", 
        UrlDepositeOrderSberPay: "https://securepayments.sberbank.ru/payment/rest/deposit.do", 
        UrlReverseSberPay: "https://securepayments.sberbank.ru//payment/rest/refund.do",
        UrlRefundSberPay: "https://securepayments.sberbank.ru//payment/rest/refund.do",
        ConnectBD: connectDb, 
       },
    }
}