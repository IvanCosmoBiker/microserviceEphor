package automatevent

import (
	 "encoding/json"
   _ "orm/model"
   "database/sql"
   "strconv"
)
var nameTable string = "automat_event"
var nameSchema string = "account"
var accountId int = 0

type AutomatEventModel struct {
    Id int
    Automat_id sql.NullInt32
    Operator_id sql.NullInt32
    Date sql.NullString
    Modem_date sql.NullInt32
    Fiscal_date sql.NullString
    Update_date sql.NullString
    Type        sql.NullInt32
    Category    sql.NullInt32
    Select_id   sql.NullString
    Ware_id     sql.NullInt32
    Name sql.NullString
    Payment_device sql.NullString
    Price_list  sql.NullInt32
    Value sql.NullInt32
    Credit sql.NullInt32
    Tax_system  sql.NullInt32
    Tax_rate sql.NullInt32
    Tax_value sql.NullInt32
    Fn        sql.NullInt64
    Fd        sql.NullInt32
    Fp        sql.NullInt64
    Fp_string sql.NullString
    Id_fr     sql.NullString
    Status sql.NullInt32
    Point_id    sql.NullInt32
    Loyality_type sql.NullInt32
    Loyality_code sql.NullString
    Error_detail  sql.NullString
    Warehouse_id  sql.NullString
    Type_fr       sql.NullInt32  
}

func (aem *AutomatEventModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &aem)
}

func (aem *AutomatEventModel) JsonSerialize() []byte{
     data, _ := json.Marshal(aem)
     return data
}

func (aem *AutomatEventModel) GetNameTable() string{
    return nameTable
}

func (aem *AutomatEventModel) GetNameSchema() string {
    accountString := strconv.Itoa(accountId)
    return nameSchema+""+accountString
}