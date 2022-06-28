package automatevent

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreAutomatEvent struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model AutomatEventModel

func (ae StoreAutomatEvent) PrepareInsertValue(value, Field interface{}) string {
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

func (ae StoreAutomatEvent) PrepareInsert(parametrs map[string]interface{}) string {
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
        if i == 1 {
            if exist {
                result += fmt.Sprintf("%s",ae.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",ae.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (ae StoreAutomatEvent) PrepareValue(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int,float64,float32,int8,int16,int32,int64:
		result = fmt.Sprintf("%q=%v",Field,value)
		case string:
		result = fmt.Sprintf("%q='%v'",Field,value)
		case nil:
		result = fmt.Sprintf("%q='%v'",Field,"null")
	}
	return result
}

func (ae StoreAutomatEvent) PrepareFieldInsertSql() string {
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

func (ae StoreAutomatEvent) PrepareFieldSql() string {
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

func (ae StoreAutomatEvent) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := ae.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (ae StoreAutomatEvent) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := ae.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (ae StoreAutomatEvent) GetDataOfMap(sql string,model bool) ([]map[string]interface{},[]AutomatEventModel) {
    ctx := context.Background()
    rows,err := ae.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    ModelAutomatEvents := []AutomatEventModel{}
    var result []map[string]interface{}
    for rows.Next(){
        ModelAutomatEventNew := AutomatEventModel{}
        err := rows.Scan(&ModelAutomatEventNew.Id, 
            &ModelAutomatEventNew.Automat_id, 
            &ModelAutomatEventNew.Operator_id, 
            &ModelAutomatEventNew.Date, 
            &ModelAutomatEventNew.Modem_date, 
            &ModelAutomatEventNew.Fiscal_date, 
            &ModelAutomatEventNew.Update_date,
            &ModelAutomatEventNew.Type,        
            &ModelAutomatEventNew.Category,    
            &ModelAutomatEventNew.Select_id,   
            &ModelAutomatEventNew.Ware_id,     
            &ModelAutomatEventNew.Name, 
            &ModelAutomatEventNew.Payment_device, 
            &ModelAutomatEventNew.Price_list,  
            &ModelAutomatEventNew.Value, 
            &ModelAutomatEventNew.Credit, 
            &ModelAutomatEventNew.Tax_system,  
            &ModelAutomatEventNew.Tax_rate, 
            &ModelAutomatEventNew.Tax_value, 
            &ModelAutomatEventNew.Fn,
            &ModelAutomatEventNew.Fd ,       
            &ModelAutomatEventNew.Fp,
            &ModelAutomatEventNew.Fp_string, 
            &ModelAutomatEventNew.Id_fr,     
            &ModelAutomatEventNew.Status, 
            &ModelAutomatEventNew.Point_id,    
            &ModelAutomatEventNew.Loyality_type, 
            &ModelAutomatEventNew.Loyality_code, 
            &ModelAutomatEventNew.Error_detail,  
            &ModelAutomatEventNew.Warehouse_id,  
            &ModelAutomatEventNew.Type_fr)
        if err != nil{
            fmt.Println(err)
            continue
        }
        ModelAutomatEvents = append(ModelAutomatEvents,ModelAutomatEventNew)
    }
    if model == true {
        return nil,ModelAutomatEvents
    }
    AutomatEventMap := make(map[string]interface{})
    for _,entry := range ModelAutomatEvents{
        AutomatEventMap["id"] = entry.Id
		AutomatEventMap["automat_id"],_ = entry.Automat_id.Value()
        AutomatEventMap["operator_id"],_ = entry.Operator_id.Value()
        AutomatEventMap["date"],_ = entry.Date.Value()
        AutomatEventMap["modem_date"],_ = entry.Modem_date.Value()
        AutomatEventMap["fiscal_date"],_ = entry.Fiscal_date.Value()
        AutomatEventMap["update_date"],_ = entry.Update_date.Value()
        AutomatEventMap["type"],_ = entry.Type.Value()
        AutomatEventMap["category"],_ = entry.Category.Value()
        AutomatEventMap["select_id"],_ = entry.Select_id.Value()
        AutomatEventMap["ware_id"],_ = entry.Ware_id.Value()
        AutomatEventMap["name"],_ = entry.Name.Value()
        AutomatEventMap["payment_device"],_ = entry.Payment_device.Value()
        AutomatEventMap["price_list"],_ = entry.Price_list.Value()
        AutomatEventMap["value"],_ = entry.Value.Value()
        AutomatEventMap["credit"],_ = entry.Credit.Value()
        AutomatEventMap["tax_system"],_ = entry.Tax_system.Value()
        AutomatEventMap["tax_rate"],_ = entry.Tax_rate.Value()
        AutomatEventMap["tax_value"],_ = entry.Tax_value.Value()
        AutomatEventMap["fn"],_ = entry.Fn.Value()
        AutomatEventMap["fd"],_ = entry.Fd.Value()
        AutomatEventMap["fp"],_ = entry.Fp.Value()
        AutomatEventMap["fp_string"],_ = entry.Fp_string.Value()
        AutomatEventMap["id_fr"],_ = entry.Id_fr.Value()
        AutomatEventMap["status"],_ = entry.Status.Value()
        AutomatEventMap["point_id"],_ = entry.Point_id.Value()
        AutomatEventMap["loyality_type"],_ = entry.Loyality_type.Value()
        AutomatEventMap["loyality_code"],_ = entry.Loyality_code.Value()
        AutomatEventMap["error_detail"],_ = entry.Error_detail.Value()
        AutomatEventMap["warehouse_id"],_ = entry.Warehouse_id.Value()
        AutomatEventMap["type_fr"],_ = entry.Type_fr.Value()
        result = append(result,AutomatEventMap)
    }
    return result,nil
}

func (ae StoreAutomatEvent) Get() ([]AutomatEventModel){
    // ctx := context.Background()
    // fieldSql := ae.PrepareFieldSql()
    commandModel := []AutomatEventModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, ae.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (ae StoreAutomatEvent) Set(parametrs map[string]interface{}){

}

func (ae StoreAutomatEvent) GetWithOptions(options map[string]interface{})([]AutomatEventModel){
    where := ae.PrepareWhere(options) 
    fieldSql := ae.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    _,modelResult := ae.GetDataOfMap(sql,true)
    return modelResult
}

func (ae StoreAutomatEvent) SetWithOptions(options map[string]interface{}){

}

func (ae StoreAutomatEvent) AddByParams(parametrs map[string]interface{}){
    accountId = parametrs["account_id"].(int)
    log.Println(accountId)
    ctx := context.Background()
    fieldSql := ae.PrepareFieldInsertSql()
    insertSql := ae.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  ae.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (ae StoreAutomatEvent) SetByParams(parametrs map[string]interface{}){
   accountId = parametrs["account_id"].(int)
   Where := ""
   options := make(map[string]interface{})
   update := ae.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = ae.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,Where)
       ae.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (ae StoreAutomatEvent) GetOneById(id interface{}) {

}