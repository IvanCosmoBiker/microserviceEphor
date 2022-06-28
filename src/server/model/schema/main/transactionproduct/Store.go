package transactionproduct

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreTransactionProduct struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model TransactionProductModel

func (tp *StoreTransactionProduct) PrepareInsertValue(value, Field interface{}) string {
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

func (tp *StoreTransactionProduct) PrepareInsert(parametrs map[string]interface{}) string {
    result := "("
    v := reflect.ValueOf(Model)
    fmt.Println(parametrs)
    typeOfS := v.Type()
    for i := 0; i< v.NumField(); i++ {
        if typeOfS.Field(i).Name == "Id"{
            continue;
        }
        field := typeOfS.Field(i).Name
        lowField := fmt.Sprintf("%v",strings.ToLower(field))
        value,exist := parametrs[lowField]
        fmt.Printf("%T",parametrs[lowField])
        if i == 1 {
            if exist {
                result += fmt.Sprintf("%s",tp.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",tp.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (tp *StoreTransactionProduct) PrepareValue(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int:
		result = fmt.Sprintf("%q=%v",Field,value)
		case string:
		result = fmt.Sprintf("%q='%v'",Field,value)
		case nil:
		result = fmt.Sprintf("%q=%v",Field,"null")
	}
	return result
}

func (tp *StoreTransactionProduct) PrepareFieldInsertSql() string {
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

func (tp *StoreTransactionProduct) PrepareFieldSql() string {
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

func (tp *StoreTransactionProduct) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := tp.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (tp *StoreTransactionProduct) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := tp.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (tp *StoreTransactionProduct) Get() ([]TransactionProductModel){
    // ctx := context.Background()
    // fieldSql := tp.PrepareFieldSql()
    commandModel := []TransactionProductModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, tp.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (tp *StoreTransactionProduct) Set(parametrs map[string]interface{}){

}

func (tp *StoreTransactionProduct) GetWithOptions(options map[string]interface{})([]TransactionProductModel){
    ctx := context.Background()
    where := tp.PrepareWhere(options) 
    fieldSql := tp.PrepareFieldSql()
    log.Println(fieldSql)
    transactionProducts := []TransactionProductModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    rows,err := tp.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    for rows.Next(){
        transactionProduct := TransactionProductModel{}
        err := rows.Scan(&transactionProduct.Id,
        &transactionProduct.Transaction_id,
        &transactionProduct.Name,
        &transactionProduct.Select_id,
        &transactionProduct.Ware_id,
        &transactionProduct.Value,
        &transactionProduct.Tax_rate,
        &transactionProduct.Quantity)
        if err != nil{
            fmt.Println(err)
            continue
        }
        transactionProducts = append(transactionProducts,transactionProduct)
    }
    return transactionProducts
}

func (tp *StoreTransactionProduct) SetWithOptions(options map[string]interface{}){

}

func (tp *StoreTransactionProduct) AddByParams(parametrs map[string]interface{}){
    ctx := context.Background()
    fieldSql := tp.PrepareFieldInsertSql()
    insertSql := tp.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  tp.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (tp *StoreTransactionProduct) SetByParams(parametrs map[string]interface{}){
   Where := ""
   options := make(map[string]interface{})
   update := tp.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = tp.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,Where)
       tp.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (tp *StoreTransactionProduct) GetOneById(id interface{}) {

}
