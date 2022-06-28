package command

import(
	"encoding/json"
    "fmt"
    "log"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    requestCommand "ephorservices/src/data/command"
    CommandRepositary "ephorservices/src/server/model/schema/main/command"
)

func SendCommandToDeviceRequest(json_data []byte) {
    CommandRequest := requestCommand.CommandServerRequest{}
    response  := make(map[string]interface{})
    json.Unmarshal(json_data, &CommandRequest)
    ConnectDb.AddLog(fmt.Sprintf("%+v",CommandRequest),"CommandSend" ,fmt.Sprintf("%+v",CommandRequest),"CommandSend")
    entry := make(map[string]interface{})
    entry["id"] = CommandRequest.Id
    entry["sended"] = CommandRepositary.SendSuccess
    log.Println(CommandRequest.Id)
    command := CommandStore.GetOneById(CommandRequest.Id)
    CommandStore.SetByParams(entry)
    cmd,_ := command.Command.Value()
    sum,_ := command.Command_param1.Value()
    response["a"] = cmd.(int64)
    response["m"]  = 2
    response["sum"] = sum.(int64)
    data, _ := json.Marshal(response)
    Rabbit.PublishMessage(data,fmt.Sprintf("ephor.1.dev.%v",CommandRequest.Imei))
}

var  Rabbit *ConnectionRabbitMQ.ChannelMQ
var  CommandStore  CommandRepositary.StoreCommand
var  ConnectDb connectionPostgresql.DatabaseInstance

func Init(rabbit *ConnectionRabbitMQ.ChannelMQ, connectPg connectionPostgresql.DatabaseInstance) {
    ConnectDb = connectPg
    Rabbit = rabbit
    CommandStore.Connection = connectPg
}
