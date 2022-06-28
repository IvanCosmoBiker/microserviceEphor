package payApi

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "strconv"
    "log"
    config "ephorservices/src/configs"
    transportHttp "ephorservices/src/pkg/transportprotocols/http/v1"
    requestApi "ephorservices/src/data/requestApi"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    transactionLogic "ephorservices/src/server/middleLayer/transaction"
)

func handler(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
        case "POST":
            json_data, _ := ioutil.ReadAll(req.Body)
            defer req.Body.Close()
            bankChannel := make(chan bool)
            requestCheck := requestApi.Request{}
            json.Unmarshal(json_data, &requestCheck)
            connectDb.AddLog(fmt.Sprintf("%+v",requestCheck),"TransactionSystem" ,fmt.Sprintf("%+v",requestCheck),"EphorErp")
            resultFoundActiveTid := transactionLogic.CheckTransactionOfAutomat(requestCheck)
            response := make(map[string]interface{})
            if resultFoundActiveTid == true {
                response["message"] = "Действие над автоматом производится другим пользователем, пожалуйста подождите"
                body, err := json.Marshal(response)
                if err != nil {
                    return
                }
                w.Write(body)
                return 
            }else {
                noise,Id := transactionLogic.AddTransaction(requestCheck)
                requestCheck.IdTransaction = strconv.Itoa(Id)
                fmt.Println(requestCheck.IdTransaction)
                response["message"] = "ok"
                response["tid"] = noise
                body, err := json.Marshal(response)
                if err != nil {
                    return
                }
                w.Write(body)
                go transactionLogic.StartTransaction(bankChannel, requestCheck)
                return 
            }
            return 
        case "GET":
            fmt.Fprintf(w, "%s: Running\n", "Servirce Payment")
            log.Println("Running")
        default:
            fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
    }
}

var (
     Config     *config.Config
     Transport *transportHttp.ServerHttp
     connectDb  connectionPostgresql.DatabaseInstance
)


func InitUrl(conf *config.Config, transport *transportHttp.ServerHttp,connectPg connectionPostgresql.DatabaseInstance) {
    Config = conf
    connectDb = connectPg
    Transport = transport
    route :=  Transport.SetHandlerListener("/pay",handler)
    route.Methods("GET","POST")
}