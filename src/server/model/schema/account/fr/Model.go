package fr

import (
   "encoding/json"
   "database/sql"
   "strconv"
)
var nameTable string = "fr"
var nameSchema string = "account"

type FrModel struct {
    Id int
    Name	sql.NullString
    Type    sql.NullInt32	
    Dev_interface	sql.NullInt32	 
    Login		sql.NullString	
    Password	sql.NullString		 
    Phone		sql.NullString		 
    Email		sql.NullString		
    Dev_addr	sql.NullString		 
    Dev_port	sql.NullInt32		
    Ofd_addr	sql.NullString		 
    Ofd_port	sql.NullInt32		 
    Inn		sql.NullString		
    Auth_public_key sql.NullString	
    Auth_private_key sql.NullString	 
    Sign_private_key sql.NullString	 
    Param1	sql.NullString		 
    Use_sn sql.NullInt32			
    Add_fiscal sql.NullInt32		
    Id_shift		sql.NullString	
    Fr_disable_cash sql.NullInt32		
    Fr_disable_cashless sql.NullInt32
}

func (fm *FrModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &fm)
}

func (fm *FrModel) JsonSerialize() []byte{
     data, _ := json.Marshal(fm)
     return data
}

func (fm *FrModel) GetNameTable() string{
    return nameTable
}

func (fm *FrModel) GetNameSchema(accountId int) string {
    accountString := strconv.Itoa(accountId)
    return nameSchema+""+accountString
}

