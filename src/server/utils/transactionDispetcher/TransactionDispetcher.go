package transactionDispetcher 

import (
    "sync"
    "log"
)

type TransactionDispetcher struct {
    mutex sync.Mutex
    transactions map[int] interface{}
}

/*
    key is composed by accountId and automatId
 */
func (t *TransactionDispetcher) AddReplayProtection(key ,automat int) bool {
    t.mutex.Lock() 
    t.transactions[key] = automat
    t.mutex.Unlock()
    return true
}

func (t *TransactionDispetcher) GetReplayProtection(key int) interface{} {
    t.mutex.Lock() 
    automat,exist := t.transactions[key]
    if exist == false {
        t.mutex.Unlock()
        return false
    }
    t.mutex.Unlock()
    return automat
}

func (t *TransactionDispetcher) RemoveReplayProtection(key int) bool {
    t.mutex.Lock() 
    _,exist := t.transactions[key]
    if exist == false {
        t.mutex.Unlock()
        return false
    }
    delete(t.transactions, key);
    t.mutex.Unlock()
    return true
}

func (t *TransactionDispetcher) AddChannel(key int) (chan []byte) {
    t.mutex.Lock() 
    context := make(chan []byte)
    t.transactions[key] = context
    t.mutex.Unlock()
    return context
}

func (t *TransactionDispetcher) RemoveChannel(key int) bool {
    t.mutex.Lock() 
    channel,exist := t.transactions[key].(chan []byte)
    if exist == false {
        t.mutex.Unlock()
        return false
    }
    close(channel)
    delete(t.transactions, key)
    t.mutex.Unlock()
    return true
}

func (t *TransactionDispetcher) Send(key int,message []byte) bool {
    t.mutex.Lock() 
    channel,exist := t.transactions[key].(chan []byte)
    log.Printf("\n [x] %v",channel)
    if exist == false {
        t.mutex.Unlock()
        return false
    }
    channel <- message
    t.mutex.Unlock()
    return true
}

func (t *TransactionDispetcher) GetTransactions() (map[int] interface{}){
    return t.transactions
}

func New() *TransactionDispetcher{
    var newMutex sync.Mutex
    var newTransactions = make(map[int] interface{})
    return &TransactionDispetcher{
        mutex: newMutex,
        transactions: newTransactions,
    }
}

