package fiscalApi

import (
    "net/http"
    "fmt"
    "io/ioutil"
    "time"
    config "ephorservices/src/configs"
    transportHttp "ephorservices/src/pkg/transportprotocols/http/v1"
    fiscalmiddlelayer "ephorservices/src/server/middleLayer/fiscal"
)

func handler(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
        case "POST":
            json_data, _ := ioutil.ReadAll(req.Body)
            defer req.Body.Close()
            timeout := make(chan bool)
            fiscalmiddlelayer.StartFiscal(timeout, json_data)
            go func() {
                select {
                case <-time.After(Config.Services.EphorFiscal.ExecuteMinutes * time.Minute):
                    timeout <- true
                }
            }()
            return 
        case "GET":
            fmt.Fprintf(w, "%s: Running\n", "Servirce Fiscal")
        default:
        fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
    }

}


var Config *config.Config
var Transport *transportHttp.ServerHttp

func InitUrl(conf *config.Config, transport *transportHttp.ServerHttp) {
    Config = conf
    Transport = transport
    route :=  Transport.SetHandlerListener("/fiscal",handler)
    route.Methods("GET","POST")
}