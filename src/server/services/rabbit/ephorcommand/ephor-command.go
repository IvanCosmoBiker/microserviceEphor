package ephorcommand

import (
   "encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
	connectionPostgresql "ephorservices/src/pkg/connectDb"
    counter "ephorservices/src/pkg/counter"
	configEphor "ephorservices/src/configs"
	ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    requestCommand "ephorservices/src/data/command"
    commandModel "ephorservices/src/server/model/schema/main/command"
    modem "ephorservices/src/server/model/schema/main/modem"
)

func ProcessingRequest(requestCommand requestCommand.CommandRequest){
    response  := make(map[string]interface{})
    request := make(map[string]interface{})
    request["imei"] = requestCommand.D
    if requestCommand.A == 1 {
       modemArray := ModemStore.GetWithOptions(request)
       modem := modemArray[0]
       requestCommand := make(map[string]interface{})
       requestCommand["modem_id"] = modem.Id
	   requestCommand["sended"] = commandModel.SendUnSuccess
       commands := CommandStore.GetWithOptions(requestCommand)
       if len(commands) != 1 {
			command := commands[0]
			cmd,_ := command.Command.Value()
			sum,_ := command.Command_param1.Value()
			response["a"] = cmd
			response["m"] = 2
			response["sum"] = sum
			parametrs := make(map[string]interface{})
			parametrs["sended"] = 1
			CommandStore.Set(parametrs)
			data, _ := json.Marshal(response)
			ConnectionRabbit.PublishMessage(data,fmt.Sprintf("ephor.1.dev.%v",request["imei"]))
			ConnDb.SetLog(fmt.Sprintf("%+v",response))
       }
       
    }
}


func initCommandCron(forever chan bool) {
	stringQueue := cfg.Services.EphorCommand.NameQueue
	go func() {
		select {
		case <-forever:
		    ConnectionRabbit.CloseChannel(stringQueue)
			if counterGo.N == 0 {
				return 
			}
		}
	}()
	msg, _ := ConnectionRabbit.RabbitMQConsume(stringQueue)
	counterGo.Add()
	for d := range msg {
		req := requestCommand.CommandRequest{}
		log.Printf("\n [x] %s", d.Body)
		dataLog := fmt.Sprintf("%s", d.Body)
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			errData, _ := fmt.Println(err)
			log.Println(errData)
			ConnDb.AddLog(dataLog, "EphorCommand", fmt.Sprintf("%s", err),"EphorCommandError")
			continue
		}
		ConnDb.AddLog(dataLog,"EphorCommand", "DataOk",req.D)
		ProcessingRequest(req)
	}
	counterGo.Sub()
	return 
}

var cfg *configEphor.Config
var ConnectionRabbit *ConnectionRabbitMQ.ChannelMQ
var ConnDb connectionPostgresql.DatabaseInstance
var counterGo counter.Counter
var CommandStore commandModel.StoreCommand
var ModemStore modem.StoreModem
var forever = make(chan bool)
var Transactions *transactionDispetcher.TransactionDispetcher

func Start(conf *configEphor.Config,Rabbit *ConnectionRabbitMQ.ChannelMQ,Db connectionPostgresql.DatabaseInstance,transactions *transactionDispetcher.TransactionDispetcher) {
	fmt.Println("Start EphorCommand...")
	ConnectionRabbit = Rabbit 
	ConnDb = Db
	cfg = conf
    CommandStore.Connection = Db
    ModemStore.Connection = Db
	Transactions = transactions
	start(forever)
}

func ReconnectQueue() {
	initCommandCron(forever)
}

func start(forever chan bool) {
	initCommandCron(forever)
}

func stop(forever chan bool) {
	ConnectionRabbit.CloseConnectRabbit()
	if counterGo.N == 0 {
		ConnDb.CloseConnectionDb()
	} else {
		go func() {
			select {
			case <-time.After(10 * time.Second):
				forever <- true
			}
		}()
		ConnDb.CloseConnectionDb()
	}
	os.Exit(3)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}



