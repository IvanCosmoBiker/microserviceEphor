package command

import (
    "encoding/json"
   _ "orm/model"
    "database/sql"
)

const (
    SendUnSuccess = 0
    SendSuccess = 1
)

var nameTable string = "modem_command"
var nameSchema string = "main"
type CommandModel struct {
    Id int
    Modem_id sql.NullInt32
    Command sql.NullInt32
    Command_param1 sql.NullInt32
    Date sql.NullString
    Sended sql.NullInt32
}

func (cm *CommandModel) InitData(jsonData []byte){
    json.Unmarshal(jsonData, &cm)
}

func (cm *CommandModel) JsonSerialize() []byte {
     data, _ := json.Marshal(cm)
     return data
}

func (cm *CommandModel) GetNameTable() string{
    return nameTable
}

func (cm *CommandModel) GetNameSchema() string {
    return nameSchema
}