package transactionproduct

import (
	 "encoding/json"
  // _ "orm/model"
   "database/sql"
)
var nameTable string = "transaction_product"
var nameSchema string = "main"

type TransactionProductModel struct {
    Id int
    Transaction_id sql.NullInt32
    Name sql.NullString
    Select_id sql.NullString
    Ware_id sql.NullInt32
    Value sql.NullInt32
	Tax_rate sql.NullInt32
	Quantity sql.NullInt32
}

func (cm *TransactionProductModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &cm)
}

func (cm *TransactionProductModel) JsonSerialize() []byte{
     data, _ := json.Marshal(cm)
     return data
}

func (cm *TransactionProductModel) GetNameTable() string{
    return nameTable
}

func (cm *TransactionProductModel) GetNameSchema() string {
    return nameSchema
}