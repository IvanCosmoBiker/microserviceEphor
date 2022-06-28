package listener

import(
    "context"
    "net/http"
    "log"
    "time"
    "github.com/gorilla/mux"
    "fmt"
)

type ServerHttp struct {
    RouterHttp *mux.Router
    Server *http.Server
}

func (server *ServerHttp) Init(url,port string) {
    server.InitRouter()
    server.InitServer(url,port)
}

func (server *ServerHttp) StartListener() {
    if err := server.Server.ListenAndServe(); err != nil {log.Println(err)}
    log.Println("Start listeners Server")
}

func (server *ServerHttp) InitRouter() {
    server.RouterHttp = mux.NewRouter()
}

func (server *ServerHttp) InitServer(url,port string) {
    Address := fmt.Sprintf(":%s",port)
    log.Println(Address) 
    s := &http.Server{
        Addr:           Address,
        Handler:        server.RouterHttp,
        ReadTimeout:    60 * time.Second,
        WriteTimeout:   60 * time.Second,
    }
    server.Server = s
}

func (server *ServerHttp) SetHandlerListener(address string, handler func(w http.ResponseWriter, req *http.Request)) *mux.Route {
    router := server.RouterHttp.HandleFunc(address, handler)
    return router
}
/*
    Завершает работу http сервера аккуратно. Закрывает соединения для новых и ждёт завершения текущих соединений используя пакет context. 
*/
func (server *ServerHttp) CloseListener() {
    ctx := context.Background()
    if err := server.Server.Shutdown(ctx); err != nil {
        // Error from closing listeners, or context timeout:
        log.Printf("HTTP server Shutdown: %v", err)
    }else {
        log.Print("HTTP server Shutdown completed successfully")
    }
}

