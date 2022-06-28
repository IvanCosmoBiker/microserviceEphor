package automatlocation 

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreAutomatLocation struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model AutomatLocationModel


func (al *StoreAutomatLocation) PrepareInsertValue(value, Field interface{}) string {
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

func (al *StoreAutomatLocation) PrepareInsert(parametrs map[string]interface{}) string {
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
                result += fmt.Sprintf("%s",al.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",al.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (al *StoreAutomatLocation) PrepareValue(value, Field interface{}) string {
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

func (al *StoreAutomatLocation) PrepareFieldInsertSql() string {
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

func (al *StoreAutomatLocation) PrepareFieldSql() string {
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

func (al *StoreAutomatLocation) PrepareValueWhere(value, Field interface{}) string {
    result := ""
	switch value.(type) {
		case int,float64,float32,int8,int16,int32,int64:
		result = fmt.Sprintf("%q=%v",Field,value)
		case string:
		result = fmt.Sprintf("%q='%v'",Field,value)
		case nil:
		result = fmt.Sprintf("%q %s",Field,"IS NULL")
	}
	return result
}

func (al *StoreAutomatLocation) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := al.PrepareValueWhere(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (al *StoreAutomatLocation) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := al.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}
   		
func (al *StoreAutomatLocation) GetDataOfMap(sql string) ([]map[string]interface{}) {
    ctx := context.Background()
    rows,err := al.Connection.Conn.Query(ctx,sql)
    if err != nil {
        log.Println(err) 
    }
    defer rows.Close()
    AutomatLocations := []AutomatLocationModel{}
    var result []map[string]interface{}
    for rows.Next(){
        AutomatLocation := AutomatLocationModel{}
        err := rows.Scan(&AutomatLocation.Id, 
            &AutomatLocation.Automat_id,
            &AutomatLocation.Company_point_id,
            &AutomatLocation.From_date,  
            &AutomatLocation.To_date)
        if err != nil{
            fmt.Println(err)
            continue
        }
        AutomatLocations = append(AutomatLocations,AutomatLocation)
    }
    automatLocationMap := make(map[string]interface{})
    for _,entry := range AutomatLocations{
        automatLocationMap["id"] = entry.Id
		automatLocationMap["automat_id"],_ = entry.Automat_id.Value()
        automatLocationMap["company_point_id"],_ = entry.Company_point_id.Value()
        automatLocationMap["from_date"],_ = entry.From_date.Value()
        automatLocationMap["to_date"],_ = entry.To_date.Value()
        result = append(result,automatLocationMap)
    }
    return result
}

func (al *StoreAutomatLocation) Get() ([]AutomatLocationModel){
    // ctx := context.Background()
    // fieldSql := al.PrepareFieldSql()
    commandModel := []AutomatLocationModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, al.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (al *StoreAutomatLocation) Set(parametrs map[string]interface{}){

}

func (al *StoreAutomatLocation) GetWithOptions(options map[string]interface{})([]map[string]interface{}){
    var result []map[string]interface{}
    accountId := options["account_id"].(int)
    delete(options, "account_id")
    where := al.PrepareWhere(options) 
    fieldSql := al.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId),Model.GetNameTable(),where)
    log.Println(sql)
    result = al.GetDataOfMap(sql)
    return result
}

func (al *StoreAutomatLocation) SetWithOptions(options map[string]interface{}){

}

func (al *StoreAutomatLocation) AddByParams(parametrs map[string]interface{}){
    accountId := parametrs["account_id"].(int)
    delete(parametrs, "account_id")
    ctx := context.Background()
    fieldSql := al.PrepareFieldInsertSql()
    insertSql := al.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(accountId),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  al.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (al *StoreAutomatLocation) SetByParams(parametrs map[string]interface{}){
   accountId := parametrs["account_id"].(int)
   delete(parametrs, "account_id")
   Where := ""
   options := make(map[string]interface{})
   update := al.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = al.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(accountId),Model.GetNameTable(),update,Where)
       al.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (al *StoreAutomatLocation) GetOneById(id,accountId interface{}) ([]map[string]interface{},bool) {
    idAutomat := id.(int)
    fieldSql := al.PrepareFieldSql()
    Where := fmt.Sprintf(" WHERE id = %v",idAutomat)
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId.(int)),Model.GetNameTable(),Where)
    result := al.GetDataOfMap(sql)
    return result,true
}