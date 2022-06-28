package command

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "log"
    "strings"
    "fmt"
    "context"
)

type StoreCommand struct {
    Connection connectionPostgresql.DatabaseInstance
    Data []CommandModel
}

var Model CommandModel

func (sc *StoreCommand) PrepareValue(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int,float64,float32,int8,int16,int32,int64:
		result = fmt.Sprintf("%q=%v",fmt.Sprintf("%v",Field),value)
		case string:
		result = fmt.Sprintf("%q='%v'",fmt.Sprintf("%v",Field),value)
		case nil:
		result = fmt.Sprintf("%q=%v",fmt.Sprintf("%v",Field),"null")
        default:
        result = fmt.Sprintf("%q='%v'",fmt.Sprintf("%v",Field),value)
	}
	return result
}

func (sc *StoreCommand) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := sc.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (sc *StoreCommand) PrepareFieldSql() string {
    result := ""
    v := reflect.ValueOf(Model)
    typeOfS := v.Type()
    for i := 0; i< v.NumField(); i++ {
        if result == "" {
            field := fmt.Sprintf("%q ", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }else {
            field := fmt.Sprintf(",%q ", typeOfS.Field(i).Name)
            result += strings.ToLower(field)
        }  
    }
    return result
}

func (sc *StoreCommand) PrepareWhere(options map[string]interface{}) string {
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

func (sc *StoreCommand) Get() ([]CommandModel) {
    ctx := context.Background()
    fieldSql := sc.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    rows,err := sc.Connection.Conn.Query(ctx,sql,nil)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    commandModel := []CommandModel{}
     
    for rows.Next(){
        comm := CommandModel{}
        err := rows.Scan(&comm.Id, &comm.Modem_id, &comm.Command, &comm.Command_param1,&comm.Date,&comm.Sended)
        if err != nil{
            fmt.Println(err)
            continue
        }
        commandModel = append(commandModel, comm)
    }

    return commandModel
}

func (sc *StoreCommand) Set(parametrs map[string]interface{}){
   commands := sc.Data
   options := make(map[string]interface{})
   update := sc.PrepareUpdate(parametrs)
   ctx := context.Background()
    for i := 0; i< len(commands); i++ {
        command := commands[i]
        options["id"] = command.Id
        where := sc.PrepareWhere(options) 
        sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(),Model.GetNameTable(),update,where)
        sc.Connection.Conn.Exec(ctx,sql)
    }
}

func (sc *StoreCommand) GetWithOptions(options map[string]interface{}) ([]CommandModel){
   ctx := context.Background()
   where := sc.PrepareWhere(options) 
   fieldSql := sc.PrepareFieldSql()
   sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
   rows,err := sc.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    commandModel := []CommandModel{}
     
    for rows.Next(){
        comm := CommandModel{}
        err := rows.Scan(&comm.Id, &comm.Modem_id, &comm.Command, &comm.Command_param1,&comm.Date,&comm.Sended)
        if err != nil{
            fmt.Println(err)
            continue
        }
        commandModel = append(commandModel, comm)
    }
    sc.Data = commandModel
    return commandModel
}

func (sc *StoreCommand) SetWithOptions(options map[string]interface{}){

}

func (sc *StoreCommand) AddByParams(parametrs map[string]interface{}){

}

func (sc *StoreCommand) SetByParams(parametrs map[string]interface{}){
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

func (sc *StoreCommand) GetOneById(id interface{}) CommandModel {
    option := make(map[string]interface{})
    option["id"] = id
    ctx := context.Background()
    fieldSql := sc.PrepareFieldSql()
    where := sc.PrepareWhere(option)
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(),Model.GetNameTable(),where)
    log.Println(sql)
    rows,err := sc.Connection.Conn.Query(ctx,sql)
    if err != nil {
            log.Println(err)     
    }
    defer rows.Close()
    commandModel := CommandModel{}
     
    for rows.Next(){
        err := rows.Scan(&commandModel.Id, &commandModel.Modem_id, &commandModel.Command, &commandModel.Command_param1,&commandModel.Date,&commandModel.Sended)
        if err != nil{
            fmt.Println(err)
            continue
        }
    }
    return commandModel
}



