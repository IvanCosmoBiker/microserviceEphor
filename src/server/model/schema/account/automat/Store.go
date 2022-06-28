package automat 

import (
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    "reflect"
    "strings"
    "context"
    "fmt"
    "log"
)

type StoreAutomat struct {
    Connection connectionPostgresql.DatabaseInstance
}

var Model AutomatModel


func (au StoreAutomat) PrepareInsertValue(value, Field interface{}) string {
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

func (au StoreAutomat) PrepareInsert(parametrs map[string]interface{}) string {
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
                result += fmt.Sprintf("%s",au.PrepareInsertValue(value,lowField))
            }else {
                result += "null"
            }
        }else {
            if exist {
                result += fmt.Sprintf(",%s",au.PrepareInsertValue(value,lowField))
            }else {
                result += ",null"
            }
        }
    }
    result += ")"
    return result
}

func (au StoreAutomat) PrepareValue(value, Field interface{}) string {
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

func (au StoreAutomat) PrepareFieldInsertSql() string {
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

func (au StoreAutomat) PrepareFieldSql() string {
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

func (au StoreAutomat) PrepareWhere(options map[string]interface{}) string {
     where := "" 
     for field, value := range options {
         ValuePrepare := au.PrepareValue(value,field)
         if where == ""{
             where += fmt.Sprintf(" WHERE %s",ValuePrepare)
         }else {
             where += fmt.Sprintf(" AND %s",ValuePrepare)
         }
     }
     return where
}

func (au StoreAutomat) PrepareUpdate(parametrs map[string]interface{}) string {
    result := ""
    for field, value := range parametrs {
        if field == "id"{
            continue
        }else {
            ValuePrepare := au.PrepareValue(value,field)
            if result == ""{
                result += fmt.Sprintf("%s ",ValuePrepare)
            }else {
                result += fmt.Sprintf(",%s ",ValuePrepare)
            }
        }
     }
    return result
}
   		
func (au StoreAutomat) GetDataOfMap(sql string) ([]map[string]interface{}) {
    ctx := context.Background()
    rows,err := au.Connection.Conn.Query(ctx,sql)
    if err != nil {
        log.Println(err) 
    }
    defer rows.Close()
    Automats := []AutomatModel{}
    var result []map[string]interface{}
    for rows.Next(){
        Automat := AutomatModel{}
        err := rows.Scan(&Automat.Id, 
            &Automat.Automat_model_id,
            &Automat.Modem_id,
            &Automat.Fr_id,  
            &Automat.Pay_system,
            &Automat.Pay_type,
            &Automat.Sbp_id, 
            &Automat.Ext1,			
            &Automat.Serial_number,	
            &Automat.Key,			 
            &Automat.Production_date,
            &Automat.From_date,	 
            &Automat.To_date,	
            &Automat.Update_date,	
            &Automat.Last_sale, 
            &Automat.Last_audit,	
            &Automat.Last_encash,	
            &Automat.Type_nosale,	
            &Automat.Type_service, 
            &Automat.Type_encashment,	  
            &Automat.Now_cash_val, 
            &Automat.Now_cashless_val, 
            &Automat.Tube_val,	
            &Automat.Now_tube_val,	
            &Automat.Control_billvalidator,	
            &Automat.Control_coinchanger,		  
            &Automat.Control_cashless,		
            &Automat.Last_coin,		
            &Automat.Last_bill,		
            &Automat.Last_cashless,	
            &Automat.Load_date,		
            &Automat.Update_config_id,	
            &Automat.Now_cash_num, 
            &Automat.Now_cashless_num, 
            &Automat.Cash_error,		
            &Automat.Cashless_error,	
            &Automat.Qr,		
            &Automat.Qr_type,	
            &Automat.Ext2,	
            &Automat.Usb1,	
            &Automat.Internet,	
            &Automat.Ethernet_mac, 		
            &Automat.Ethernet_ip,	 	
            &Automat.Ethernet_netmask,	
            &Automat.Ethernet_gateway,	
            &Automat.Faceid_type,		
            &Automat.Faceid_id,		
            &Automat.Faceid_addr,		
            &Automat.Faceid_port,		
            &Automat.Faceid_username,	
            &Automat.Faceid_password,	
            &Automat.Summ_max_fr,		
            &Automat.Last_login, 
            &Automat.Warehouse_id,	
            &Automat.Screen_text)
        if err != nil{
            fmt.Println(err)
            continue
        }
        Automats = append(Automats,Automat)
    }
    
    for _,entry := range Automats{
        automatMap := make(map[string]interface{})
        automatMap["id"] = entry.Id
		automatMap["automat_model_id"],_ = entry.Automat_model_id.Value()
        automatMap["modem_id"],_ = entry.Modem_id.Value()
        automatMap["fr_id"],_ = entry.Fr_id.Value()
        automatMap["pay_system"],_ = entry.Pay_system.Value()
        automatMap["pay_type"],_ = entry.Pay_type.Value()
        automatMap["sbp_id"],_ = entry.Sbp_id.Value()
        automatMap["ext1"],_ = entry.Ext1.Value()
        automatMap["serial_number"],_ = entry.Serial_number.Value()
        automatMap["key"],_ = entry.Key.Value()
        automatMap["production_date"],_ = entry.Production_date.Value()
        automatMap["from_date"],_ = entry.From_date.Value()
        automatMap["to_date"],_ = entry.To_date.Value()
        automatMap["update_date"],_ = entry.Update_date.Value()
        automatMap["last_sale"],_ = entry.Last_sale.Value()
        automatMap["last_audit"],_ = entry.Last_audit.Value()
        automatMap["last_encash"],_ = entry.Last_encash.Value()
        automatMap["type_nosale"],_ = entry.Type_nosale.Value()
        automatMap["type_service"],_ = entry.Type_service.Value()
        automatMap["type_encashment"],_ = entry.Type_encashment.Value()
        automatMap["now_cash_val"],_ = entry.Now_cash_val.Value()
        automatMap["now_cashless_val"],_ = entry.Now_cashless_val.Value()
        automatMap["tube_val"],_ = entry.Tube_val.Value()
        automatMap["now_tube_val"],_ = entry.Now_tube_val.Value()
        automatMap["control_billvalidator"],_ = entry.Control_billvalidator.Value()
        automatMap["control_coinchanger"],_ = entry.Control_coinchanger.Value()
        automatMap["control_cashless"],_ = entry.Control_cashless.Value()
        automatMap["last_coin"],_ = entry.Last_coin.Value()
        automatMap["last_bill"],_ = entry.Last_bill.Value()
        automatMap["last_cashless"],_ = entry.Last_cashless.Value()
        automatMap["load_date"],_ = entry.Load_date.Value()
        automatMap["update_config_id"],_ = entry.Update_config_id.Value()
        automatMap["now_cash_num"],_ = entry.Now_cash_num.Value()
        automatMap["now_cashless_num"],_ = entry.Now_cashless_num.Value()
        automatMap["cash_error"],_ = entry.Cash_error.Value()
        automatMap["cashless_error"],_ = entry.Cashless_error.Value()
        automatMap["qr"],_ = entry.Qr.Value()
        automatMap["qr_type"],_ = entry.Qr_type.Value()
        automatMap["ext2"],_ = entry.Ext2.Value()
        automatMap["usb1"],_ = entry.Usb1.Value()
        automatMap["internet"],_ = entry.Internet.Value()
        automatMap["ethernet_mac"],_ = entry.Ethernet_mac.Value()
        automatMap["ethernet_ip"],_ = entry.Ethernet_ip.Value()
        automatMap["ethernet_netmask"],_ = entry.Ethernet_netmask.Value()
        automatMap["ethernet_gateway"],_ = entry.Ethernet_gateway.Value()
        automatMap["faceid_type"],_ = entry.Faceid_type.Value()
        automatMap["faceid_id"],_ = entry.Faceid_id.Value()
        automatMap["faceid_addr"],_ = entry.Faceid_addr.Value()
        automatMap["faceid_port"],_ = entry.Faceid_port.Value()
        automatMap["faceid_username"],_ = entry.Faceid_username.Value()
        automatMap["faceid_password"],_ = entry.Faceid_password.Value()
        automatMap["summ_max_fr"],_ = entry.Summ_max_fr.Value()
        automatMap["last_login"],_ = entry.Last_login.Value()
        automatMap["warehouse_id"],_ = entry.Warehouse_id.Value()
        automatMap["screen_text"],_ = entry.Screen_text.Value()
        result = append(result,automatMap)
    }
    return result
}

func (au StoreAutomat) Get() ([]AutomatModel){
    // ctx := context.Background()
    // fieldSql := au.PrepareFieldSql()
    commandModel := []AutomatModel{}
    // sql := fmt.Sprintf("SELECT %v FROM %v.%v",fieldSql,Model.GetNameSchema(),Model.GetNameTable())
    // pgxtpan.Select(ctx, au.Connection.Conn, &commandModel, sql)
    return commandModel
}

func (au StoreAutomat) Set(parametrs map[string]interface{}){

}

func (au StoreAutomat) GetWithOptions(options map[string]interface{})([]map[string]interface{}){
    var result []map[string]interface{}
    accountId := options["account_id"].(int)
    delete(options, "account_id")
    where := au.PrepareWhere(options) 
    fieldSql := au.PrepareFieldSql()
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId),Model.GetNameTable(),where)
    result = au.GetDataOfMap(sql)
    return result
}

func (au StoreAutomat) SetWithOptions(options map[string]interface{}){

}

func (au StoreAutomat) AddByParams(parametrs map[string]interface{}){
    accountId := parametrs["account_id"].(int)
    delete(parametrs, "account_id")
    ctx := context.Background()
    fieldSql := au.PrepareFieldInsertSql()
    insertSql := au.PrepareInsert(parametrs)
    sql := fmt.Sprintf("INSERT INTO %v.%v (%v) VALUES %s ",Model.GetNameSchema(accountId),Model.GetNameTable(),fieldSql,insertSql)
    log.Println(sql)
    commandTag, err :=  au.Connection.Conn.Exec(ctx,sql)
    if err != nil {
        fmt.Println(err)
    }
    if commandTag.RowsAffected() != 1 {
        fmt.Println(commandTag.RowsAffected())
    }
}

func (au StoreAutomat) SetByParams(parametrs map[string]interface{}){
   accountId := parametrs["account_id"].(int)
   delete(parametrs, "account_id")
   Where := ""
   options := make(map[string]interface{})
   update := au.PrepareUpdate(parametrs)
   ctx := context.Background()
   _,exist := parametrs["id"]
   if exist {
       options["id"] = parametrs["id"]
       Where = au.PrepareWhere(options)
       sql := fmt.Sprintf("UPDATE %v.%v SET %s %v",Model.GetNameSchema(accountId),Model.GetNameTable(),update,Where)
       au.Connection.Conn.Exec(ctx,sql)
   }else {
       return 
   }
}

func (au StoreAutomat) GetOneById(id,accountId interface{}) ([]map[string]interface{},bool) {
    idAutomat := id.(int)
    fieldSql := au.PrepareFieldSql()
    Where := fmt.Sprintf(" WHERE id = %v",idAutomat)
    sql := fmt.Sprintf("SELECT %v FROM %v.%v %v",fieldSql,Model.GetNameSchema(accountId.(int)),Model.GetNameTable(),Where)
    result := au.GetDataOfMap(sql)
    return result,true
}