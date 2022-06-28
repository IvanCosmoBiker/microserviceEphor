package automat 

import (
   "encoding/json"
   "database/sql"
   "strconv"
)
var nameTable string = "automat"
var nameSchema string = "account"

type AutomatModel struct {
    Id int
    Automat_model_id sql.NullInt32
    Modem_id sql.NullInt32
    Fr_id   sql.NullInt32
    Pay_system sql.NullInt32
    Pay_type sql.NullInt32
    Sbp_id  sql.NullInt32
    Ext1	sql.NullInt32		
    Serial_number 	sql.NullString
    Key		sql.NullString		 
    Production_date	sql.NullString 
    From_date	sql.NullString 
    To_date		sql.NullString
    Update_date	sql.NullString
    Last_sale sql.NullString
    Last_audit	sql.NullString
    Last_encash	sql.NullString
    Type_nosale	sql.NullInt32
    Type_service sql.NullInt32
    Type_encashment	 sql.NullInt32 
    Now_cash_val sql.NullInt32
    Now_cashless_val sql.NullInt32
    Tube_val	sql.NullInt64
    Now_tube_val	sql.NullInt64
    Control_billvalidator	sql.NullInt32
    Control_coinchanger		sql.NullInt32  
    Control_cashless		sql.NullInt32
    Last_coin		sql.NullString
    Last_bill		sql.NullString
    Last_cashless	sql.NullString
    Load_date		sql.NullString
    Update_config_id	sql.NullString
    Now_cash_num sql.NullInt32
	Now_cashless_num sql.NullInt32
    Cash_error		sql.NullInt32
    Cashless_error	sql.NullInt32
    Qr		sql.NullInt32
    Qr_type	sql.NullInt32
    Ext2	sql.NullInt32
    Usb1	sql.NullInt32
    Internet	sql.NullInt32
    Ethernet_mac sql.NullString		
    Ethernet_ip	 sql.NullString	
    Ethernet_netmask	sql.NullString
    Ethernet_gateway	sql.NullString
    Faceid_type		sql.NullInt32
    Faceid_id		sql.NullString
    Faceid_addr		sql.NullString
    Faceid_port		sql.NullInt32
    Faceid_username	sql.NullString
    Faceid_password	sql.NullString
    Summ_max_fr		sql.NullInt64
    Last_login sql.NullString
    Warehouse_id	sql.NullString
    Screen_text		sql.NullString
}

func (am *AutomatModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &am)
}

func (am *AutomatModel) JsonSerialize() []byte{
     data, _ := json.Marshal(am)
     return data
}

func (am *AutomatModel) GetNameTable() string{
    return nameTable
}

func (am AutomatModel) GetNameSchema(accountId int) string {
    accountString := strconv.Itoa(accountId)
    return nameSchema+""+accountString
}