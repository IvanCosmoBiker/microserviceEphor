package transaction

import (
   "encoding/json"
   "database/sql"
)
var nameTable string = "transaction"
var nameSchema string = "main"

type TransactionModel struct {
    Id int
    Noise sql.NullInt32
    Token_id sql.NullString
    Token_type sql.NullInt32
    Account_id sql.NullString
    Automat_id sql.NullString
    Date sql.NullString
    Sum sql.NullInt32
    Ps_type sql.NullInt32
	Ps_order sql.NullString
	Ps_code sql.NullString
    Ps_desc sql.NullString
    Ps_invoice_id sql.NullString
    Pay_type sql.NullInt32
    Fn sql.NullInt32
    Fd sql.NullInt32
    Fp sql.NullString
    F_type sql.NullInt32
    F_receipt sql.NullString
    F_desc sql.NullString
    F_status sql.NullInt32
    Qr_format sql.NullInt32
    F_qr sql.NullString
    Status sql.NullInt32
    Error sql.NullString
}

func (tm *TransactionModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &tm)
}

func (tm *TransactionModel) JsonSerialize() []byte{
     data, _ := json.Marshal(tm)
     return data
}

func (tm *TransactionModel) GetNameTable() string{
    return nameTable
}

func (tm *TransactionModel) GetNameSchema() string {
    return nameSchema
}