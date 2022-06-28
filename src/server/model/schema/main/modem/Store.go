package modem

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "github.com/georgysavva/scany/pgxscan"
    "context"
    "fmt"
    "log"
)

type StoreModem struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model ModemModel

func (sc *StoreModem) PrepareValue(value, Field interface{}) string {
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

func (sc *StoreModem) PrepareFieldSql() string {
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

func (sc *StoreModem) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := sc.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (sc *StoreModem) Get() ([]ModemModel){
    ctx := context.Background()
    fieldSql := sc.PrepareFieldSql()
    commandModel := []ModemModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    pgxscan.Select(ctx, sc.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (sc *StoreModem) Set(parametrs map[string]interface{}){

}

func (sc *StoreModem) GetWithOptions(options map[string]interface{})([]ModemModel){
    ctx := context.Background()
    where := sc.PrepareWhere(options) 
    fieldSql := sc.PrepareFieldSql()
    log.Println(fieldSql)
    commandModel := []ModemModel{}
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    rows,err := sc.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    for rows.Next(){
        comm := ModemModel{}
        err := rows.Scan(&comm.Id,&comm.Account_id,&comm.Imei,&comm.Hash,&comm.Nonce,&comm.Hardware_version,&comm.Software_version,&comm.Software_release,&comm.Phone,&comm.Signal_quality,&comm.Last_login,&comm.Last_ex_id,&comm.Ipaddr,&comm.Static,&comm.Gsm_apn,&comm.Gsm_username,&comm.Gsm_password,&comm.Dns1,&comm.Dns2,&comm.Add_date)
        if err != nil{
            fmt.Println(err)
            continue
        }
        commandModel = append(commandModel,comm)
    }
    return commandModel
}

func (sc *StoreModem) SetWithOptions(options map[string]interface{}){

}

func (sc *StoreModem) AddByParams(parametrs map[string]interface{}){

}

func (sc *StoreModem) SetByParams(parametrs map[string]interface{}){

}

func (sc *StoreModem) GetOneById(id interface{}) {

}
