package rabbitmq

import(
	"fmt"
	"log"
	"amqp"
	"time"
)

type Config struct {
	Login string
	Password string
	Address string
	Port string
	MaxAttempts int
}

type ChannelMQ struct {
	Channel         *amqp.Channel
	ConnectionChannel ConnectionRabbit
	ChannelConn bool
	ConfigRabbit  Config
}

type ConnectionRabbit struct {
	Connect  *amqp.Connection
	IsReady         bool
}

func (ch *ChannelMQ) SetConfig(login string, password string, address string, port string,maxAttempts int){
	ch.ConfigRabbit.Login = login
	ch.ConfigRabbit.Password = password
	ch.ConfigRabbit.Address = address
	ch.ConfigRabbit.Port = port
	ch.ConfigRabbit.MaxAttempts = maxAttempts
}

func (ch *ChannelMQ) ConnectionToRabbit(login string, password string, address string, port string,maxAttempts int) error {
	stringConnection:= fmt.Sprintf("amqp://%s:%s@%s:%s",login,password,address,port)
	conn, err := amqp.Dial(stringConnection)
	ch.SetConfig(login,password,address,port,maxAttempts)
	if err != nil {
		return err
	}else {
		ch.ConnectionChannel.Connect = conn
	    ch.ConnectionChannel.IsReady = true
		return nil
	}
	return nil
}

func (ch *ChannelMQ) Reconnect(rabbitreconnect chan bool) {
	for { 
		result := ch.CheckConnect(rabbitreconnect)
		if result == true {
			rabbitreconnect <- result
		}
		time.Sleep(10 * time.Second)
	}
}

func (ch *ChannelMQ) CheckConnect(rabbitreconnect chan bool) bool { 
	var err error
	if ch.ConnectionChannel.Connect.IsClosed() == true {
		if err = ch.ConnectionToRabbit(ch.ConfigRabbit.Login,ch.ConfigRabbit.Password,ch.ConfigRabbit.Address,ch.ConfigRabbit.Port,ch.ConfigRabbit.MaxAttempts); err != nil {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				log.Println("reconnecting...")
				err = ch.ConnectionToRabbit(ch.ConfigRabbit.Login,ch.ConfigRabbit.Password,ch.ConfigRabbit.Address,ch.ConfigRabbit.Port,ch.ConfigRabbit.MaxAttempts)
				if err == nil {
					ch.ConnectQueue()
					return true
				}
				log.Printf("connection was lost. Error: %s. Waiting for 10 sec...\n", err)
			}
		}else {
			return false
		}
	}else {
		return false
	}
	return false
}

func (ch *ChannelMQ) QueueDeclareRabbit(nameQueue string) {
	_,err := ch.Channel.QueueDeclare(
		nameQueue, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		ch.ConnectQueue()
		ch.QueueDeclareRabbit(nameQueue)
		return 
	}else {
		failOnError(err, "Failed to register a queue")
		return 
	}
}

func (ch *ChannelMQ) ConnectQueue(){
	channel, err := ch.ConnectionChannel.Connect.Channel()
	if err != nil {
		ch.ChannelConn = false
	}
	failOnError(err, "Failed to open a channel")
	ch.Channel = channel
	ch.ChannelConn = true
}
 
// func (ch *ChannelMQ) GetDeliver() amqp.Delivery {
// 	return amqp.Delivery
// }

func (ch *ChannelMQ)  NotifyReturn() (chan amqp.Return){
	channelToReturn := make(chan amqp.Return)
	result := ch.Channel.NotifyReturn(channelToReturn)
	return result
}	

func (ch *ChannelMQ) RabbitMQConsume(nameQueue string) (<-chan amqp.Delivery, error){
	
	msgs, err := ch.Channel.Consume(
		nameQueue, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	  if err != nil {
		ch.QueueDeclareRabbit(nameQueue)
		ch.RabbitMQConsume(nameQueue)
	    ch.Channel.QueueBind(
			nameQueue, // queue name
			nameQueue,     // routing key
			"amq.topic", // exchange
			false,
			nil,
		)
		return nil,err
	  }else {
		ch.Channel.QueueBind(
			nameQueue, // queue name
			nameQueue,     // routing key
			"amq.topic", // exchange
			false,
			nil,
		)
		return msgs,err
	  }
	
}

func (ch *ChannelMQ) PublishMessage(body []byte,name string) error {
	err := ch.Channel.Publish(
	"amq.topic",     // exchange
	name, // routing key
	true,  // mandatory
	false,  // immediate
	amqp.Publishing {
		ContentType: "application/json",
		Body:        body,
		DeliveryMode: amqp.Persistent,
	})
	return err
}

func (ch *ChannelMQ) CloseChannel(name string){
	ch.Channel.Cancel(name,true)
}

func (ch *ChannelMQ) CloseConnectRabbit() {
	ch.Channel.Close()
	ch.ConnectionChannel.Connect.Close()
	ch.ChannelConn = false
	ch.ConnectionChannel.IsReady = false
}

func failOnError(err error, msg string) {
	if err != nil {
	  log.Printf("%s: %s", msg, err)
	}
}
