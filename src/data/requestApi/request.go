package requestApi

import (
    "encoding/json"
)

type Request struct {
    Config struct {
        UserPhone string
        ReturnUrl string 
        DeepLink string
        Login string
		Password string
        TokenType int
        BankType int
        PayType int
        CurrensyCode int
        Language string
        Description string
        AccountId int
        AutomatId int
        DeviceType int
        QrFormat int
    }
    Products []map[string]interface{}
    Date string
    OrderId string
    SbolBankInvoiceId string
    IdTransaction string
    MerchantId string
    GateWay string
    PaymentToken string
    WareId string
    Sum int
    SumOneProduct int
    SumMax int
    Imei string
}

type Data struct {
	Events             []string
	Imei, CheckId, Inn string
	ConfigFR           struct {
		Host, Cert, Key, Sign, Port,Login,Password string
		Fiscalization int
	}
    TypeFr int
	InQueue int
	Fields  struct {
		Request json.RawMessage
	}
    DataResponse map[string]interface{}
}