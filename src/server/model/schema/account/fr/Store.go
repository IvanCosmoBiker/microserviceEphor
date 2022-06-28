package fr 

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreFr struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model FrModel


func (fr *StoreFr) PrepareInsertValue(value, Field interface{}) string {
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

func (fr *StoreFr) PrepareInsert(parametrs map[string]interface{}) string {
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
                result += fmt.Sprintf("%s",fr.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",fr.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (fr *StoreFr) PrepareValue(value, Field interface{}) string {
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

func (fr *StoreFr) PrepareFieldInsertSql() string {
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

func (fr *StoreFr) PrepareFieldSql() string {
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

func (fr *StoreFr) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := fr.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (fr *StoreFr) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := fr.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}

func (fr *StoreFr) GetDataOfMap(sql string) ([]map[string]interface{}) {
    ctx := context.Background()
    rows,err := fr.Connection.Conn.Query(ctx,sql)
    if err != nil {
        log.Println(err) 
    }
    defer rows.Close()
    Frs := []FrModel{}
    var result []map[string]interface{}
    for rows.Next(){
        Fr := FrModel{}
        err := rows.Scan(&Fr.Id, 
            &Fr.Name,
            &Fr.Type,
            &Fr.Dev_interface,  
            &Fr.Login,
            &Fr.Password,
            &Fr.Phone, 
            &Fr.Email,			
            &Fr.Dev_addr,	
            &Fr.Dev_port,			 
            &Fr.Ofd_addr,
            &Fr.Ofd_port,	 
            &Fr.Inn,	
            &Fr.Auth_public_key,	
            &Fr.Auth_private_key, 
            &Fr.Sign_private_key,	
            &Fr.Param1,	
            &Fr.Use_sn,	
            &Fr.Add_fiscal, 
            &Fr.Id_shift,	  
            &Fr.Fr_disable_cash, 
            &Fr.Fr_disable_cashless)
        if err != nil{
            fmt.Println(err)
            continue
        }
        Frs = append(Frs,Fr)
    }
    FrMap := make(map[string]interface{})
    for _,entry := range Frs{
        FrMap["id"] = entry.Id
		FrMap["name"],_ = entry.Name.Value()
        FrMap["type"],_ = entry.Type.Value()
        FrMap["dev_interface"],_ = entry.Dev_interface.Value()
        FrMap["login"],_ = entry.Login.Value()
        FrMap["password"],_ = entry.Password.Value()
        FrMap["phone"],_ = entry.Phone.Value()
        FrMap["email"],_ = entry.Email.Value()
        FrMap["dev_addr"],_ = entry.Dev_addr.Value()
        FrMap["dev_port"],_ = entry.Dev_port.Value()
        FrMap["ofd_addr"],_ = entry.Ofd_addr.Value()
        FrMap["ofd_port"],_ = entry.Ofd_port.Value()
        FrMap["inn"],_ = entry.Inn.Value()
        FrMap["auth_public_key"],_ = entry.Auth_public_key.Value()
        FrMap["auth_private_key"],_ = entry.Auth_private_key.Value()
        FrMap["sign_private_key"],_ = entry.Sign_private_key.Value()
        FrMap["param1"],_ = entry.Param1.Value()
        FrMap["use_sn"],_ = entry.Use_sn.Value()
        FrMap["add_fiscal"],_ = entry.Add_fiscal.Value()
        FrMap["id_shift"],_ = entry.Id_shift.Value()
        FrMap["fr_disable_cash"],_ = entry.Fr_disable_cash.Value()
        FrMap["fr_disable_cashless"],_ = entry.Fr_disable_cashless.Value()
        result = append(result,FrMap)
    }
    log.Printf("%s",result)
    return result
}

func (fr *StoreFr) Get() ([]FrModel){
    // ctx := context.Background()
    // fieldSql := fr.PrepareFieldSql()
    commandModel := []FrModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, fr.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (fr *StoreFr) Set(parametrs map[string]interface{}){

}

func (fr *StoreFr) GetWithOptions(options map[string]interface{})([]map[string]interface{}){
    var result []map[string]interface{}
    accountId := options["account_id"].(int)
    delete(options, "account_id")
    where := fr.PrepareWhere(options) 
    fieldSql := fr.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId),Model.GetNameTable(),where)
    result = fr.GetDataOfMap(sql)
    return result
}

func (fr *StoreFr) SetWithOptions(options map[string]interface{}){

}

func (fr *StoreFr) AddByParams(parametrs map[string]interface{}){
    accountId := parametrs["account_id"].(int)
    delete(parametrs, "account_id")
    ctx := context.Background()
    fieldSql := fr.PrepareFieldInsertSql()
    insertSql := fr.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(accountId),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  fr.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (fr *StoreFr) SetByParams(parametrs map[string]interface{}){
   accountId := parametrs["account_id"].(int)
   delete(parametrs, "account_id")
   Where := ""
   options := make(map[string]interface{})
   update := fr.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = fr.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(accountId),Model.GetNameTable(),update,Where)
       fr.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (fr *StoreFr) GetOneById(id interface{},accountId interface{}) ([]map[string]interface{},bool) {
    options := make(map[string]interface{})
    options["id"] = id
    fieldSql := fr.PrepareFieldSql()
    Where := fr.PrepareWhere(options)
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId.(int)),Model.GetNameTable(),Where)
    log.Printf("%s",sql)
    result := fr.GetDataOfMap(sql)
    return result,true
}