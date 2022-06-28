package ware

import (
	"encoding/json"
   _ "orm/model"
   "database/sql"
   "strconv"
)
var nameTable string = "ware"
var nameSchema string = "account"
var accountId int = 0

type WareModel struct {
    Id int
    Code  sql.NullString
    Name     sql.NullString
    State    sql.NullInt32
    Type     sql.NullInt32
    Description  sql.NullString
}

func (wm *WareModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &wm)
}

func (wm *WareModel) JsonSerialize() []byte{
     data, _ := json.Marshal(wm)
     return data
}

func (wm *WareModel) GetNameTable() string{
    return nameTable
}

func (wm *WareModel) GetNameSchema() string {
    accountString := strconv.Itoa(accountId)
    return nameSchema+""+accountString
}