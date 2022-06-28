package transaction 

import ()

var (
    Prepayment = 1 // предоплата товара 
    Postpaid = 2 // постоплата товара
)

var (
    IceboxStatus_Drink = 1 // выдача напитков
    IceboxStatus_Icebox = 2 // дверь открыта
    IceboxStatus_End = 3 // выдача завершена
)

var (
    VendState_Session	 		= 1 //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
    VendState_Approving   		= 2 //[14] Продукт выбран. Ожидание оплаты.
    VendState_Vending	 		= 3 //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
    VendState_VendOk	     	= 4 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
    VendState_VendError     	= 5 //[13]PAY_OK_AUTOMAT_PREPARED Оплата успешна, ТА приготовил продукт
)

var (
    VendError_VendFailed         = 769 //769 Ошибка выдачи продукта
    VendError_SessionCancelled   = 770 //770
    VendError_SessionTimeout     = 771 //771
    VendError_WrongProduct       = 772 //772
    VendError_VendCancelled      = 773 //773
    VendError_ApprovingTimeout   = 774 //774
)

var (
    TransactionState_Idle 				= 0; // Transaction Idle
    TransactionState_MoneyHoldStart 	= 1; // создали транзакцию банка
    TransactionState_MoneyHoldWait 		= 2; // ожидает ответ от банка
    TransactionState_VendSessionStart 	= 3; // PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
    TransactionState_VendSession	 	= 4; //[11]PAY_OK_BUTTON_PRESS Оплата успешна, ожидание нажатия пользователем кнопки на ТА
    TransactionState_VendApproving   	= 5; //[14] Продукт выбран. Ожидание оплаты.
    TransactionState_Vending	 		= 6; //[12]PAY_OK_AUTOMAT_PREPARE Оплата успешна, ТА готовит продукт
    TransactionState_MoneyDebitStart	= 8;
    TransactionState_MoneyDebitWait		= 9;
    TransactionState_MoneyDebitOk		= 10;
    TransactionState_VendOk				= 11; // приготовил продукт
	TransactionState_MoneyDebit 		= 12;
    TransactionState_Error 				= 120;
    TransactionState_ErrorTimeOut       = 121;
	TransactionState_WaitFiscal			= 14
)
var (
	 TypeTokenGooglePay 		= 1
	 TypeTokenApplePay 			= 2
	 TypeTokenSamsungPay 		= 3
	 TypeTokenSberPayWeb 		= 4
	 TypeTokenSberPayAndroid 	= 5
	 TypeTokenSberPayiOS 		= 6
)

type Transaction struct {
    Tid string
	UserPhone string
	ReturnUrl string 
	DeepLink string
	TokenType int
	Noise int
	Date string
    Products []map[string]interface{}
    Sum int
    PayType int // Prepayment or Postpaid
    DeviceType int
    AccountId int
    AutomatId int
    Token string
	Tax_system int
    SumMax int
	Address string
	PointName string
	QrFormat int
	PointId int
}

func (t Transaction) GetDescriptionCodeCooler(code int) string {
	stringCode := ""
	 switch code {
		 case IceboxStatus_Drink:
		 stringCode = `Подойдите к кофе машине. Подставьте стакан и выберите напиток`
		 return stringCode
		 fallthrough
		 case IceboxStatus_Icebox:
		 stringCode = `Замок открыт, заберите продукты`
		 return stringCode
		 fallthrough
		 case IceboxStatus_End:
		 stringCode = `Замок закрыт. Счастливого пути`
		 return stringCode
	 }
	 return stringCode
}

func (t Transaction) GetStatusServerCooler(status int) int {
	 switch status {
		 case IceboxStatus_Drink:
		 return TransactionState_VendSession
		 fallthrough
		 case IceboxStatus_Icebox:
		 return TransactionState_VendSession
		 fallthrough
		 case IceboxStatus_End:
		 return TransactionState_VendOk
	 }
	 return TransactionState_Error
}

func (t Transaction) GetDescriptionCode(code int) string {
	stringCode := ""
	 switch code {
		 case VendState_Session:
		 stringCode = `Оплата успешна, ожидание нажатия пользователем кнопки на ТА`
		 return stringCode 
		 fallthrough
		 case VendState_Approving:
		 stringCode = `Продукт выбран. Ожидание оплаты.`
		 return stringCode
		 fallthrough
		 case VendState_Vending:
		 stringCode = `Оплата успешна, ТА готовит продукт`
		 return stringCode
		 fallthrough
		 case VendState_VendOk:
		 stringCode = `Оплата успешна, ТА приготовил продукт`
		 return stringCode
		 fallthrough
		 case VendState_VendError:
		 stringCode = `Ошибка`
		 return stringCode
         fallthrough
         case TransactionState_ErrorTimeOut:
		 stringCode = `Время ответа от автомата истекло`
		 return stringCode
		 fallthrough
		 case IceboxStatus_Drink:
		 stringCode = `Подойдите к кофе машине. Подставьте стакан и выберите напиток`
		 return stringCode
	 }
	 return stringCode
}

func (t Transaction) GetStatusServer(status int) int {
	 switch status {
		 case VendState_Session:
		 return TransactionState_VendSession
		 fallthrough
		 case VendState_Approving:
		 return TransactionState_VendApproving
		 fallthrough
		 case VendState_Vending:
		 return TransactionState_Vending
		 fallthrough
		 case VendState_VendOk:
		 return TransactionState_VendOk
		 fallthrough
		 case VendState_VendError:
		 return TransactionState_Error
	 }
	 return TransactionState_Error
}

func (t Transaction) GetDescriptionErr(err int) string{
	stringErr := ""
	 switch err {
		 case VendError_VendFailed:
		 stringErr = `Ошибка выдачи товара на автомате`
		 return stringErr
		 fallthrough
		 case VendError_SessionCancelled:
		 stringErr = `Продажа отменена автоматом`
		 return stringErr
		 fallthrough
		 case VendError_SessionTimeout:
		 stringErr = `Время ожидание выбора товара на автомате истекло`
		 return stringErr
		 fallthrough
		 case VendError_WrongProduct:
		 stringErr = `Выбранный на автомате товар не совпадает с оплаченым`
		 return stringErr
		 fallthrough
		 case VendError_VendCancelled:
		 stringErr = `Выдача товара отменена автоматом`
		 return stringErr
		 fallthrough
		 case VendError_ApprovingTimeout:
		 stringErr = `Время ожидание оплаты истекло`
		 return stringErr
         fallthrough
         case TransactionState_ErrorTimeOut:
		 stringErr = `Время ответа от автомата истекло`
		 return stringErr
	 }
	 return stringErr
}