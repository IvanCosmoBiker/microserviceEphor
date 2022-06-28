package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
    "time"
    "log"
)

/* Configuration structure, load from config.json, global */
type Config struct {
    RabbitMq struct {
        Login, Password, Address, Port string
        MaxAttempts int
        ExecuteTimeSeconds time.Duration
    }
    Db struct {
       Login,Password,Address,DatabaseName string
       Port  uint16
       PgConnectionPool int 
    }
    Services struct {
        Http struct {
            Address,Port string
        }
        EphorPay struct {
            NameQueue string
            Bank struct {
                ExecuteMinutes time.Duration // this parametr for time run work with bank
                PollingTime time.Duration
            }
        }
        EphorCommand struct {
            NameQueue string
            ExecuteMinutes time.Duration
            Listener struct {
                ExecuteMinutes time.Duration
            }
        }
        EphorFiscal struct {
            NameQueue string
            ResponseUrl string
            ExecuteMinutes,SleepMilliSec time.Duration
            Listener struct {
                ExecuteMinutes time.Duration
            }
        }
    }
    ErrorCount int
    LogFile string
    ExecuteMinutes time.Duration // this parametr work execute time for one transaction
    Debug bool
}

func (c *Config) Load() {
	file, _ := os.Open("config.json")
	byteValue, _ := ioutil.ReadAll(file)
	defer file.Close()
	json.Unmarshal(byteValue, &c)
    log.Printf("%v",c)
}