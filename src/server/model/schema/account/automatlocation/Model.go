package automatlocation

import (
   "encoding/json"
   "database/sql"
   "strconv"
)
var nameTable string = "automat_location"
var nameSchema string = "account"

type AutomatLocationModel struct {
    Id int
    Automat_id sql.NullInt32
    Company_point_id sql.NullInt32
    From_date   sql.NullString
    To_date sql.NullString
}

func (am *AutomatLocationModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &am)
}

func (am *AutomatLocationModel) JsonSerialize() []byte{
     data, _ := json.Marshal(am)
     return data
}

func (am *AutomatLocationModel) GetNameTable() string{
    return nameTable
}

func (am *AutomatLocationModel) GetNameSchema(accountId int) string {
    accountString := strconv.Itoa(accountId)
    return nameSchema+""+accountString
}