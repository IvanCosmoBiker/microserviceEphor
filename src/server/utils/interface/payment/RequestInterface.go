package payment

import (

)

type Request struct {
    Config struct {
        BankType int
        PayType int
        CurrensyCode int
        TokenType int
        Language string
        Description string
        AccountId int
        AutomatId int
        DeviceType int
    }
    Products []map[string]interface{}
    Date string
    IdTransaction string
    MerchantId string
    GateWay string
    PaymentToken string
    WareId string
    Sum int
    SumMax int
    Imei string
}