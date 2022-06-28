package ware

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreWare struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model WareModel

func (ws *StoreWare) PrepareInsertValue(value, Field interface{}) string {
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

func (ws *StoreWare) PrepareInsert(parametrs map[string]interface{}) string {
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
                result += fmt.Sprintf("%s",ws.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",ws.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (ws *StoreWare) PrepareValue(value, Field interface{}) string {
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

func (ws *StoreWare) PrepareFieldInsertSql() string {
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

func (ws *StoreWare) PrepareFieldSql() string {
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

func (ws *StoreWare) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := ws.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (ws *StoreWare) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := ws.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (ws *StoreWare) GetDataOfMap(sql string,model bool) ([]map[string]interface{},[]WareModel) {
    ctx := context.Background()
    rows,err := ws.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    ModelWares := []WareModel{}
    var result []map[string]interface{}
    for rows.Next(){
        WareModelNew := WareModel{}
        err := rows.Scan(&WareModelNew.Id,
        &WareModelNew.Code, 
        &WareModelNew.Name, 
        &WareModelNew.State, 
        &WareModelNew.Type, 
        &WareModelNew.Description)
        if err != nil{
            fmt.Println(err)
            continue
        }
        ModelWares = append(ModelWares,WareModelNew)
    }
    if model == true {
        return nil,ModelWares
    }
    WareMap := make(map[string]interface{})
    for _,entry := range ModelWares {
        WareMap["id"] = entry.Id
		WareMap["Code"],_ = entry.Code.Value()
        WareMap["Name"],_ = entry.Name.Value()
        WareMap["State"],_ = entry.State.Value()
        WareMap["Type"],_ = entry.Type.Value()
        WareMap["Description"],_ = entry.Description.Value()
        result = append(result,WareMap)
    }
    return result,nil
}

func (ws *StoreWare) Get() ([]WareModel){
    // ctx := context.Background()
    // fieldSql := ws.PrepareFieldSql()
    commandModel := []WareModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, ws.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (ws *StoreWare) Set(parametrs map[string]interface{}){

}

func (ws *StoreWare) GetWithOptions(options map[string]interface{})([]WareModel){
    accountId = options["account_id"].(int)
    delete(options, "account_id")
    where := ws.PrepareWhere(options)
    log.Println(where) 
    fieldSql := ws.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    _,modelResult := ws.GetDataOfMap(sql,true)
    return modelResult
}

func (ws *StoreWare) SetWithOptions(options map[string]interface{}){

}

func (ws *StoreWare) AddByParams(parametrs map[string]interface{}){
    accountId = parametrs["account_id"].(int)
    log.Println(accountId)
    ctx := context.Background()
    fieldSql := ws.PrepareFieldInsertSql()
    insertSql := ws.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  ws.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (ws *StoreWare) SetByParams(parametrs map[string]interface{}){
   accountId = parametrs["account_id"].(int)
   Where := ""
   options := make(map[string]interface{})
   update := ws.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = ws.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,Where)
       ws.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (ws *StoreWare) GetOneById(id interface{}) {

}