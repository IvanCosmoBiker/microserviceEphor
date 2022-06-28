package transaction

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreTransaction struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model TransactionModel

func (t *StoreTransaction) PrepareInsertValue(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int,float64,float32,int8,int16,int32,int64:
		result = fmt.Sprintf("%v",value)
		case string:
		result = fmt.Sprintf("'%v'",value)
		case nil:
		result = fmt.Sprintf("%v","null")

	}
	return result
}

func (t *StoreTransaction) PrepareInsert(parametrs map[string]interface{}) string {
    result := "("
    v := reflect.ValueOf(Model)
    typeOfS := v.Type()
    for i := 0; i< v.NumField(); i++ {
        if typeOfS.Field(i).Name == "Id"{
            continue;
        }
        field := typeOfS.Field(i).Name
        lowField := fmt.Sprintf("%v",strings.ToLower(field))
        //fmt.Println(parametrs[lowField])
        value,exist := parametrs[lowField]
        if i == 1 {
            if exist {
                result += fmt.Sprintf("%s",t.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",t.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (t *StoreTransaction) PrepareValue(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int,float64,float32,int8,int16,int32,int64:
		result = fmt.Sprintf("%q=%v",Field,value)
		case string:
		result = fmt.Sprintf("%q='%v'",Field,value)
		case nil:
		result = fmt.Sprintf("%q=%v",Field,"null")
	}
	return result
}

func (t *StoreTransaction) PrepareFieldInsertSql() string {
    result := ""
    v := reflect.ValueOf(Model)
    typeOfS := v.Type()
     for i := 0; i< v.NumField(); i++ {
        field := fmt.Sprintf("%q", typeOfS.Field(i).Name)
        if typeOfS.Field(i).Name == "Id"{
            continue;
        }
        if result == "" {
            field = fmt.Sprintf("%q", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }else {
            field = fmt.Sprintf(",%q ", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }  
    }
    return result
}

func (t *StoreTransaction) PrepareFieldSql() string {
    result := ""
    v := reflect.ValueOf(Model)
    typeOfS := v.Type()
     for i := 0; i< v.NumField(); i++ {
        if result == "" {
            field := fmt.Sprintf("%q", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }else {
            field := fmt.Sprintf(",%q ", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }  
    }
    return result
}

func (t *StoreTransaction) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := t.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (t *StoreTransaction) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := t.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (t *StoreTransaction) Get() ([]TransactionModel){
    // ctx := context.Background()
    // fieldSql := t.PrepareFieldSql()
    commandModel := []TransactionModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, t.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (t *StoreTransaction) Set(parametrs map[string]interface{}){

}

func (t *StoreTransaction) GetWithOptions(options map[string]interface{})([]TransactionModel){
    ctx := context.Background()
    where := t.PrepareWhere(options) 
    fieldSql := t.PrepareFieldSql()
    log.Println(fieldSql)
    transactions := []TransactionModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    rows,err := t.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    for rows.Next(){
        transaction := TransactionModel{}
        err := rows.Scan(&transaction.Id,
        &transaction.Token_id,
        &transaction.Token_type,
        &transaction.Account_id,
        &transaction.Automat_id,
        &transaction.Date,
        &transaction.Sum,
        &transaction.Ps_type,
        &transaction.Ps_order,
        &transaction.Ps_code,
        &transaction.Ps_desc,
        &transaction.Ps_invoice_id,
        &transaction.Pay_type,
        &transaction.Fn,
        &transaction.Fd,
        &transaction.Fp,
        &transaction.F_type,
        &transaction.F_receipt,
        &transaction.F_desc,
        &transaction.F_status,
        &transaction.Qr_format,
        &transaction.F_qr,
        &transaction.Status,
        &transaction.Error)
        if err != nil{
            fmt.Println(err)
            continue
        }
        transactions = append(transactions,transaction)
    }
    return transactions
}

func (t *StoreTransaction) SetWithOptions(options map[string]interface{}){

}

func (t *StoreTransaction) AddByParams(parametrs map[string]interface{}) int {
    var id int
    ctx := context.Background()
    fieldSql := t.PrepareFieldInsertSql()
    insertSql := t.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s RETURNING id",Model.GetNameSchema(),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    err :=  t.Connection.Conn.QueryRow(ctx,sql).Scan(&id)
    if err != nil {
        fmt.Println(err)
    }
   return id
}

func (t *StoreTransaction) SetByParams(parametrs map[string]interface{}){
   Where := ""
   options := make(map[string]interface{})
   update := t.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = t.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,Where)
       fmt.Println(sql)
       t.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (t *StoreTransaction) GetOneById(id interface{}) {

}
