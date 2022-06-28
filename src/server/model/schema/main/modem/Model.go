package modem

import (
	 "encoding/json"
   _ "orm/model"
   "database/sql"
)
var nameTable string = "modem"
var nameSchema string = "main"
type ModemModel struct {
    Id int
    Account_id int
    Imei sql.NullString
    Hash sql.NullString
    Nonce sql.NullString
    Hardware_version sql.NullInt32
	Software_version sql.NullInt32
	Software_release sql.NullInt32
	Phone			sql.NullString		
	Signal_quality 	sql.NullInt32	 
	Last_login	 	sql.NullString	 
	Last_ex_id    	sql.NullInt32     
	Ipaddr			sql.NullString 
	Static			sql.NullInt32 
	Gsm_apn			sql.NullString	 
	Gsm_username	sql.NullString	 
	Gsm_password	sql.NullString	 
	Dns1			sql.NullString	 
	Dns2			sql.NullString	 
	Add_date		sql.NullString
}

func (cm *ModemModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &cm)
}

func (cm *ModemModel) JsonSerialize() []byte{
     data, _ := json.Marshal(cm)
     return data
}

func (cm *ModemModel) GetNameTable() string{
    return nameTable
}

func (cm *ModemModel) GetNameSchema() string {
    return nameSchema
}
    
    
    
    
