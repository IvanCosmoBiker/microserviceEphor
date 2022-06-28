package commandApi

import (
    "net/http"
    "fmt"
    "io/ioutil"
	"time"
    config "ephorservices/src/configs"
    transportHttp "ephorservices/src/pkg/transportprotocols/http/v1"
    commandmiddlelayer "ephorservices/src/server/middleLayer/command"
)

func handler(w http.ResponseWriter, req *http.Request) {
        switch req.Method {
        case "POST":
            json_data, _ := ioutil.ReadAll(req.Body)
            defer req.Body.Close()
            commandmiddlelayer.SendCommandToDeviceRequest(json_data)
            go func() {
                select {
                    case <-time.After(time.Duration(Config.Services.EphorCommand.Listener.ExecuteMinutes) * time.Minute):
                        return 
                }
            }()
            return 
        case "GET":
            fmt.Fprintf(w, "%s: Running\n", "Command")
        default:
            fmt.Fprintf(w, "Sorry, only POST and GET method is supported.")
        }
}

var Config          *config.Config
var Transport       *transportHttp.ServerHttp

func InitUrl(conf *config.Config, transport *transportHttp.ServerHttp) {
    Config = conf
    Transport = transport
    route := Transport.SetHandlerListener("/command",handler)
    route.Methods("GET","POST")
}