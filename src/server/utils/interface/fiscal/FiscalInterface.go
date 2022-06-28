package fiscal

import (
	 transactionStruct "ephorservices/src/data/transaction"
	 requestFiscal "ephorservices/src/data/requestApi"
)
var (
     Fr_None = 0
	 Fr_PayOnlineFA = 1
	 Fr_KaznachejFA = 2
	 Fr_RPSystem1FA = 3
	 Fr_TerminalFA = 4
	 Fr_OrangeData = 5
	 Fr_ChekOnline = 6
	 Fr_EphorOrangeData = 7
	 Fr_EphorOnline = 8
	 Fr_NanoKassa = 9
	 Fr_ServerOrangeData = 10
	 Fr_EphorServerOrangeData = 11
	 Fr_OFD = 12
	 Fr_ServerNanoKassa = 13
	 Fr_ProSistem = 14
	 Fr_CheckBox = 15
)
// -- Status Fiscal -- //
var (
	 Status_None		 	 = 0; // продажа за 0 рублей или касса отключена
	 Status_Complete		 = 1; // чек создан успешно, реквизиты получены
	 Status_InQueue		 = 2; // чек добавлен в очередь, реквизиты не получены
	 Status_Unknown		 = 3; // результат постановки чека в очередь не известен
	 Status_Error			 = 4; // ошибка создания чека
	 Status_Overflow		 = 5; // очередь отложенной регистрации переполнена
	 Status_Manual			 = 6; // чек фискализирован вручную
	 Status_Need			 = 7; // чек необходимо фискализировать
	 Status_MAX			 = 7;
	 Status_MAX_CHECK		 = 8;
	 Status_OFF_FR		 	 = 9; // оключение фискализации со стороны клиента
	 Status_OFF_DA		 	 = 10; // отключение фискализации для безналичной оплаты
	 Status_OFF_CA		 	 = 11; // отключение фискализации для оплаты за наличку
	 Status_Check_Cancel	 = 12; // отмена чека
    )
// -- Error Fiscal -- //
var (
	 Status_Fr_Fr_InQueue		 = 101 // чек добавлен в очередь, реквизиты не получены
	 Status_Fr_Unknown		 = 102 // результат постановки чека в очередь не известен
	 Status_Fr_Error		 = 103 // ошибка создания чека
	 Status_Fr_Overflow	 = 104 // очередь отложенной регистрации переполнена
	 Status_Fr_MAX			 = 105 // превышена максимальная сумма чека
    )
type Fiscal interface {
    InitData(transactionStruct.Transaction,map[string]interface{})
    SendCheck() map[string]interface{}
    GetStatus(orderId ...string) map[string]interface{}
	GetQrPicture(date string, summ int, frResponse map[string]interface{}) string 
	GetQrUrl(date string, summ int, frResponse map[string]interface{}) string
	GetStatusApi(data requestFiscal.Data) map[string]interface{}
	SendCheckApi(data requestFiscal.Data) (map[string]interface{},requestFiscal.Data)
}