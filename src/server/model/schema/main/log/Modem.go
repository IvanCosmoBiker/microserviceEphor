package log

import (
	 "encoding/json"
   _ "orm/model"
   "database/sql"
)

var nameTable string = "log"
var nameSchema string = "main"

type LogModel struct {
    Id              int
    Address         sql.NullString
    Login           sql.NullString
    Date            sql.NullString
    Request_id      sql.NullString
    Request_uri     sql.NullString
	Request_data    sql.NullString
	Response        sql.NullString
	Runtime			sql.NullInt32
	Runtime_details sql.NullInt32
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
    