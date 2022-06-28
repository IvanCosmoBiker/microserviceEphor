package log

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "github.com/georgysavva/scany/pgxscan"
    "context"
    "fmt"
    "log"
)

type StoreLog struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model LogModel

func (sl *StoreLog) PrepareValue(value, Field interface{}) string {
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
func (sl *StoreLog) PrepareInsertValue(value, Field interface{}) string {
     result := ""
	switch value.(type) {
		case int:
		result = fmt.Sprintf("%v",value)
		case string:
		result = fmt.Sprintf("'%v'",value)
		case nil:
		result = fmt.Sprintf("%v","null")
	}
	return result
}

func (sl *StoreLog) PrepareInsert(parametrs map[string]interface{}) string {
    result := ""
    v := reflect.ValueOf(Model)
    typeOfS := v.Type()
     for i := 0; i< v.NumField(); i++ {
        field := fmt.Sprintf("%q", typeOfS.Field(i).Name)
        lowField := strings.ToLower(field)
        _,exist := parametrs[lowField] 
        if result == "" {
            if exist {
                result += fmt.Sprintf("%s",s.PrepareInsertValue())
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",s.PrepareInsertValue())
            }else {
                result += ",null"
            }
        }
    }
    return result
}

func (sl *StoreLog) PrepareFieldSql() string {
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

func (sl *StoreLog) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := sl.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (sl *StoreLog) Get() ([]LogModel){
    ctx := context.Background()
    fieldSql := sl.PrepareFieldSql()
    logModel := []LogModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    pgxscan.Select(ctx, sl.Connection.Conn, &logModel, sql)
    return logModel
}

func (sl *StoreLog) Set(parametrs map[string]interface{}){

}

func (sl *StoreLog) GetWithOptions(options map[string]interface{})([]LogModel){
    ctx := context.Background()
    where := sl.PrepareWhere(options) 
    fieldSql := sl.PrepareFieldSql()
    log.Println(fieldSql)
    logModel := []LogModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    rows,err := sl.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    for rows.Next(){
        log := LogModel{}
        err := rows.Scan(&log.Id,&log.Address,&log.Login,&log.Date,&log.Request_id,&log.Request_uri,&log.Request_data,&log.Response,&log.Runtime,&log.Runtime_details)
        if err != nil{
            fmt.Println(err)
            continue
        }
        logModel = append(logModel,log)
    }
    return logModel
}

func (sl *StoreLog) SetWithOptions(options map[string]interface{}){

}

func (sl *StoreLog) AddByParams(parametrs map[string]interface{}){
    ctx := context.Background()
    options := make(map[string]interface{})
    fieldSql := sl.PrepareFieldSql()
    insertSql := sc.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(),Model.GetNameTable(),fieldSql,insertSql)
    sc.Connection.Conn.Exec(ctx,sql)
}

func (sl *StoreLog) SetByParams(parametrs map[string]interface{}){
    Where := ""
    options := make(map[string]interface{})
    update := sc.PrepareUpdate(parametrs)
    ctx := context.Background()
    _,exist := parametrs["id"]
    if exist {
        options["id"] = parametrs["id"]
        Where = sc.PrepareWhere(options)
        sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,Where)
        sc.Connection.Conn.Exec(ctx,sql)
    }else {
        return 
    }
}

func (sl *StoreLog) GetOneById(id interface{}) {

}