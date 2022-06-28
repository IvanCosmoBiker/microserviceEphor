package microserviceEphor

import (
    "time"
    "log"
    config "ephorservices/src/configs"
    transactionDispetcher "ephorservices/src/server/utils/transactionDispetcher"
    connectionPostgresql "ephorservices/src/pkg/connectDb"
    ConnectionRabbitMQ "ephorservices/src/pkg/rabbitmq"
    ephorpay "ephorservices/src/server/services/rabbit/ephorpay"
    ephorcommand "ephorservices/src/server/services/rabbit/ephorcommand"
    ephorfiscal "ephorservices/src/server/services/rabbit/ephorfiscal"
    transportHttp "ephorservices/src/pkg/transportprotocols/http/v1"
    commandApi "ephorservices/src/api/v1/commandApi"
    fiscalApi "ephorservices/src/api/v1/fiscalApi"
    payApi "ephorservices/src/api/v1/payApi"
    commandLogic "ephorservices/src/server/middleLayer/command"
    fiscalLogic "ephorservices/src/server/middleLayer/fiscal"
    transactionLogic "ephorservices/src/server/middleLayer/transaction"
)

type ManagerService struct {
    Config *config.Config
    Service *serviceEphor
    ConnectionDb  *connectionPostgresql.DatabaseInstance
    ConnectionRabbit *ConnectionRabbitMQ.ChannelMQ
    Transport  *transportHttp.ServerHttp
    TransactionDispetcher *transactionDispetcher.TransactionDispetcher
    connectRabbit chan bool
}

func (m *ManagerService) InitUrls() {
    commandApi.InitUrl(m.Config,m.Transport)
    fiscalApi.InitUrl(m.Config,m.Transport)
    payApi.InitUrl(m.Config,m.Transport,*m.ConnectionDb)
    m.Transport.StartListener()
}

func (m *ManagerService) InitMiddleLayer() {
    commandLogic.Init(m.ConnectionRabbit, *m.ConnectionDb)
    fiscalLogic.Init(m.ConnectionDb,m.Config)
    transactionLogic.Init(m.Config,*m.ConnectionDb,m.ConnectionRabbit,m.TransactionDispetcher)
}

func (m *ManagerService) InitListenerHttp() {
    var httpServer transportHttp.ServerHttp
    m.Transport = &httpServer
    m.Transport.Init(m.Config.Services.Http.Address,m.Config.Services.Http.Port)
}

func (m *ManagerService) ConnectionDataBase(conf *config.Config) {
    m.ConnectionDb = &connectionPostgresql.DatabaseInstance{}
	m.ConnectionDb.NewConn(conf.Db.PgConnectionPool, conf.Db.Login, conf.Db.Password, conf.Db.Address, conf.Db.DatabaseName, conf.Db.Port)
	_,err := m.ConnectionDb.GetConn()
    if err != nil {log.Printf("%s",err)}
}

func (m *ManagerService) ConnectRabbit(conf *config.Config) {
    m.ConnectionRabbit = &ConnectionRabbitMQ.ChannelMQ{}
    err := m.ConnectionRabbit.ConnectionToRabbit(conf.RabbitMq.Login, conf.RabbitMq.Password, conf.RabbitMq.Address, conf.RabbitMq.Port,conf.RabbitMq.MaxAttempts)
    if err != nil {log.Printf("%v",err)}
    m.ConnectionRabbit.ConnectQueue()
}

func (m *ManagerService) checkStatusRabbit() {
     go func(){
        for {
            select {
                case <-m.connectRabbit:
                log.Println("Reconect Queues")
                go ephorpay.ReconnectQueue()
                go ephorcommand.ReconnectQueue()
                go ephorfiscal.ReconnectQueue()
            }
            time.Sleep(10 * time.Second)
        }
    }()
}

func (m *ManagerService) InitServices(s *serviceEphor) {
    m.Service = s
    m.Config = m.Service.ConfigFile
    m.connectRabbit = make(chan bool)
    m.TransactionDispetcher = transactionDispetcher.New()
    m.ConnectionDataBase(m.Service.ConfigFile)
    m.ConnectRabbit(m.Service.ConfigFile)
    m.InitListenerHttp()
    go m.InitUrls()
    m.InitMiddleLayer()
    m.checkStatusRabbit()
    go m.ConnectionRabbit.Reconnect(m.connectRabbit)
    go ephorpay.Start(m.Service.ConfigFile,m.ConnectionRabbit,*m.ConnectionDb,m.TransactionDispetcher)
    go ephorcommand.Start(m.Service.ConfigFile,m.ConnectionRabbit,*m.ConnectionDb,m.TransactionDispetcher)
    go ephorfiscal.Start(m.Service.ConfigFile,m.ConnectionRabbit,*m.ConnectionDb,m.TransactionDispetcher)
}

func (m *ManagerService) StopService() bool {
    m.Transport.CloseListener()
    for {
        transactions := m.TransactionDispetcher.GetTransactions()
        if len(transactions) < 1 {
            break 
        }
    }
    m.ConnectionRabbit.CloseConnectRabbit()
    m.ConnectionDb.CloseConnectionDb()
    return true
}



