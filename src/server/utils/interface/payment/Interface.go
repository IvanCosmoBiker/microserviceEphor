package payment

import (
    requestApi "ephorservices/src/data/requestApi"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
)
var (
    TypeSber = 1
	TypeVendPay = 2
)

type Bank interface {
    HoldMoney(req requestApi.Request) map[string]interface{}
    DebitHoldMoney(orderId string,sum int, req requestApi.Request) map[string]interface{}
    ReturnMoney(orderId string,req requestApi.Request) map[string]interface{}
    InitBankData(connect connectionPostgresql.DatabaseInstance)
    Timeout()
}