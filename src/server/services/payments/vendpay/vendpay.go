package sber

import(
    "bytes"
	"encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
    "log"
    "math"
    interfaceBank "ephorservices/src/server/utils/interface/payment"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    requestApi "ephorservices/src/data/requestApi"
)
const TransactionState_Idle 				= 0 // Transaction Idle
const TransactionState_MoneyHoldStart 		= 1 // создали транзакцию банка
const TransactionState_MoneyHoldWait 		= 2 // ожидает ответ от банка
const TransactionState_MoneyDebitStart		= 8
const TransactionState_MoneyDebitWait		= 9
const TransactionState_MoneyDebitOk			= 10
const TransactionState_ReturnMoney          = 12
const TransactionState_Error 				= 120


const Order_Registered = 0 //- заказ зарегистрирован, но не оплачен;
const Order_HoldMoney  = 1 //- предавторизованная сумма удержана (для двухстадийных платежей);
const Order_FullAuthorizationOfTheAmount =  2 //- проведена полная авторизация суммы заказа;
const Order_AuthorizationCanceled = 3 //- авторизация отменена;
const Order_RefundOperationPerformed  = 4 //- по транзакции была проведена операция возврата;
const Order_AuthorizationThroughTheServerHasBeenInitiated = 5 //- инициирована авторизация через сервер контроля доступа банка-эмитента;
const Order_AuthorizationDenied = 6 //- авторизация отклонена.



type VendPay struct {
    Name string
    Counter int
    PaymentType int
    UrlCreateOrder string
    UrlGetStatusOrder string
    UrlCancelOrder string
    UrlReverse string
    Status int
    ConnectBD connectionPostgresql.DatabaseInstance
    Req requestApi.Request
    Res interfaceBank.Response
}

type NewVendStruct struct {
    VendPay
}

var ConnectBD connectionPostgresql.DatabaseInstance

func (v VendPay) InitBankData(connect connectionPostgresql.DatabaseInstance){
    ConnectBD = connect
    ResponsePaymentSystem["tid"] = v.Req.IdTransaction
}

func (v VendPay) HoldMoney(Request requestApi.Request) map[string]interface{} {
    sum := float64(Request.Sum)
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции",math.Round((sum/100)))
    ResponsePaymentSystem["description"] = fmt.Sprintf("Сумма %.2f удержана, ожидайте завершения транзакции",math.Round((sum/100)))
    ResponsePaymentSystem["code"] = Order_FullAuthorizationOfTheAmount
    ResponsePaymentSystem["orderId"] = "vendpay"
    return ResponsePaymentSystem
}

func (v VendPay) getStatus(Request requestApi.Request) map[string]interface{} {
    result := make(map[string]interface{})
    return result
}

func (v VendPay) DebitHoldMoney(orderId string,sum int,Request requestApi.Request) map[string]interface{} {
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Деньги списаны"
    ResponsePaymentSystem["description"] = "Деньги списаны"
    ResponsePaymentSystem["code"] = TransactionState_MoneyDebitOk
    return ResponsePaymentSystem
}

func (v VendPay) ReturnMoney(orderId string,Request requestApi.Request) map[string]interface{} {
    ResponsePaymentSystem["status"] = true
    ResponsePaymentSystem["message"] = "Деньги возвращены"
    ResponsePaymentSystem["description"] = "Деньги возвращены"
    ResponsePaymentSystem["code"] = TransactionState_ReturnMoney
    return ResponsePaymentSystem
}


func (v VendPay) Timeout(){
}

func (v VendPay) GetPaymentType(){

}

func (v VendPay) Call(method string, url string, json_request []byte) (interfaceBank.Response) {
    Response := interfaceBank.Response{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		v.Res.Code = 0
		v.Res.Status = TransactionState_Error
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), &Response)
    return Response
}
var connectDb connectionPostgresql.DatabaseInstance
var ResponsePaymentSystem map[string]interface{}

func (vend NewVendStruct) NewBank() interfaceBank.Bank  /* тип interfaceBank.Bank*/ {
    ResponsePaymentSystem = make(map[string]interface{})
    return &NewVendStruct{
        VendPay: VendPay{
        Name: "VendPay",
        Counter: 0,
        PaymentType: 1, // srandart type payment 
        Status: 0,
        UrlCreateOrder: "https://3dsec.sberbank.ru/payment/google/payment.do",
        UrlGetStatusOrder: "https://3dsec.sberbank.ru/payment/google/getOrderStatusExtended.do",
        UrlCancelOrder: "https://3dsec.sberbank.ru/payment/google/deposit.do", // use for 
        UrlReverse: "https://3dsec.sberbank.ru/payment/google/reverse.do", // use for 
        ConnectBD: connectDb,
       },
    }
}