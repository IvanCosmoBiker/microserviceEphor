package payment

import (

)

type Response struct {
    Success bool
    ActionCode int
    ErrorCode int
    ErrorMessage string
    OrderStatus  int
    Error struct {
        Code int
        Description string
        Message string
    }
    OrderId string
    BankType int
    MerchantId string
    ExternalParams map[string]interface{}
    GateWay string
    Sum int
    Status int
    Code int
    Data map[string]interface{}
}


