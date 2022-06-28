package orange

import (
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
    "bytes"
	"crypto"
	"crypto/tls"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/rand"
	"encoding/json"
	"encoding/binary"
	"encoding/pem"
	"encoding/base64"
	"math"
	"errors"
	"strings"
	"time"
	"io"
	"os"
	"strconv"
	randString "ephorservices/src/pkg/randgeneratestring"
    transactionStruct "ephorservices/src/data/transaction"
    interfaceFiscal "ephorservices/src/server/utils/interface/fiscal"
	requestFiscal "ephorservices/src/data/requestApi"
	parserTypes "ephorservices/src/pkg/parser/typeParse"
)


const (
	TaxRate_NDSNone = 0
	TaxRate_NDS0 = 1
	TaxRate_NDS10 = 2
	TaxRate_NDS18 = 3
)
const (
	TaxSystem_OSN    = 0x01 // Общая ОСН
	TaxSystem_USND   = 0x02 // Упрощенная доход
	TaxSystem_USNDMR = 0x04 // Упрощенная доход минус расход
	TaxSystem_ENVD   = 0x08 // Единый налог на вмененный доход
	TaxSystem_ESN    = 0x10 // Единый сельскохозяйственный налог
	TaxSystem_Patent = 0x20 // Патентная система налогообложения
)

var pathFileCert = "/var/www/html/test/cert"

var privateKeySign = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDXNEfOOxwuSauj
uZMsjvR0BOsr4ehjK3rbxX3wLdwh0nxAJ91vJ6nP2C54wULJCF8J+FGLKVcP07dV
h7m7o1rWXRuiLd31Flj1/gLQeZo0eJgPeACtLR5+xvL46I73vJmyHPlTKHpKvrph
9Nd7KnqkzDT6yDdGKaPuYqRn0wtW4mlyalbma6i4MEaOBb8mAmpA/PCGXUtwIj9X
CB4T5yFNWBvyexewgsoaNvUGAwrPbpvCuAflmmsdYHYIKfjv/ZO5xjudINm2adXZ
yVXMHHM0A/ACaPrgnxDbPOTG7bRutqge9L/yh1mHeHXgz9XObXlV6d1n0kmBDXEy
gdeItuDVAgMBAAECggEBALtaEXVScps9mcbkxWMSZXEn4xEGEEldzgzMp3JUioOL
eo5j5lxh3G1NGFAaeCkKN6s3Ws5bRCdMOxykF6dqdKeQ0YDki4pWVUZ7SDn007IA
lulIoNYjJJxcWaUm2WiF8gxlOw4RfD3cQ+kJvhrFBZa5DRqS+cQEdmoPyG93BTUy
N7Sp97m8D374e4mAavCjj0G316x4g3okADVi1QsPvbu4tSpx9x9iRXSzYJegLdW7
DIEStoiJGYk/mS1GjzkquH19c5hvtnRVFXXlZsdy0VL3N1UH2iEKIsDGKr7dQsbN
vP7VEVERSRT7vcFITFh4ePL/RkI/2Qrt0u9QKKQ4iNUCgYEA2kFujDBmt5GK2ZW3
w4JzKRmQVIp1+YMhqbLsiqWiKA5AwABOOuykMD5k2lhCTYN7NsFcRZsTzk0bWiYV
zj6mnjde64w3v/L2/LvXvg1oCSYAdu2oGy8lWrEHWE9JJukxZ/iSYwoVJb46SizZ
qBdhC+jTeOsI/67as7asCvqECisCgYEA/GvCSBE8A14dAz/suRR6p9cEwh8JAN8b
09GZpr4qA0SilEPOB3qUKgFBcyR0buSdD6MUKKZNcK7vqbpe88B7tCpt4G/f1xVe
NZVDsgsXbXz6BIqvrL/faZ1DOPGPjkVFMk9N8T+tfKywY+/qU2tjBKn10QXurk07
6p8OkCDyQP8CgYEAiSTGb0bWtJC63CCM+UhWTsQmgkkC+sdgdr7cjf6oV10laMCI
Z9RdE4eRXfZJq2VsHisAbSiWGHMxNcNqvk916UNH3OEeAvqMIqFyXpUUA3OipRiP
Io3MfiFxSReBEvdDOV7jtWIXicDv5b4rAsm2DIK/p2KhI/DeskCd+MQUBkMCgYEA
7velakzGn/mNRfJSzbURmav6GT0AbQ7LbXDVIgKOC6ICuJKojnQBqPKfX753bDSK
bK9a+lDWp4M16V1DX0gu1JYGh5/iLeFQ2zGAcSIG/+R9XadeQRE1FOuJJHOsEGiL
5eEmTOqX95wVMceD842KpHOzADu5htIfkzMZumE2d0kCgYB+1aC77TFg0aDQq9hM
3+lrTJ5q1DNVsq/1IZGjyAXv3Guxm/urrG75ztOWWikg+q5mXINUnBDYPu0DnJCF
chf4nYFOAGJ5VQZrVblcQxoN+X3hM1Y4/6Pzg4cPiFGU+QTQjpaajr2voS7GWt7m
cAuNQXAHU4LCkppfW+xm9g3vBw==
-----END RSA PRIVATE KEY-----`)

type Outcome struct {
	Imei string
	Data struct {
		Message, Status,Method  string
		Code, StatusCode,Fiscalization int
		Fields           struct {
			Fp, Fn string
			Fd     float64
		}
	}
}

var (
	layoutISO = "2006-01-02 15:04:05"
	layoutQr  = "20060102T150405"
)

type ResponseOrange struct {
	Code   int
	Status string
	Message string
	Errors []string
	Data   map[string]interface{}
}

type ConfigOrange struct {
	Cert string
	Host string
	Port string 
	Key string
	Group string 
	TaxSystem int
	Sign string
	AutomatNumber int
	Inn string
}

type Orange struct {
    Name string
	Config ConfigOrange
}

type NewOrangeStruct struct {
    Orange
}

func (ofd *Orange) ConvertTax(tax int) int {
	switch tax {
		 case TaxRate_NDSNone:
		 return 6
		 fallthrough
		 case TaxRate_NDS0:
		 return 5
		 fallthrough
		 case TaxRate_NDS10:
		 return 2
		 fallthrough
		 case TaxRate_NDS18:
		 return 1
	 }
	 return 6
}

func (ofd *Orange) ConvertTaxationSystem(taxsystem int) int {
	switch taxsystem {
		 case TaxSystem_OSN:
		 return 0
		 fallthrough
		 case TaxSystem_USND:
		 return 1
		 fallthrough
		 case TaxSystem_USNDMR:
		 return 2
		 fallthrough
		 case TaxSystem_ENVD:
		 return 3
		 fallthrough
		 case TaxSystem_ESN:
		 return 4
		 fallthrough
		 case TaxSystem_Patent:
		 return 5
	 }
	 return 0
}

func (ofd *Orange) MakeUrlQr(date string, summ int, frResponse map[string]interface{}) string {
	t, _ := time.Parse(layoutISO, date)
	valueSumm := summ/100
    stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1",fmt.Sprintf("%s",t.Format(layoutQr)),fmt.Sprintf("%v.00",valueSumm),frResponse["fn"],frResponse["fd"],frResponse["fp"]);
	log.Println(stringResult)
	return stringResult
}

func (ofd *Orange) GetQrPicture(date string, summ int, frResponse map[string]interface{}) string {
	result := ofd.MakeUrlQr(date,summ,frResponse)
	str := base64.StdEncoding.EncodeToString([]byte(result))
	return str
}

func (ofd *Orange) GetQrUrl(date string, summ int, frResponse map[string]interface{}) string {
	return ofd.MakeUrlQr(date,summ,frResponse)
}

func (ofd *Orange) RsaEncrypt(origData,privatesignKey []byte) ([]byte, error) {
    rng := rand.Reader
    block, _ := pem.Decode(privateKeySign)
    if block == nil {
        return nil, errors.New("is not pem")
    }
    pubInterface,err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return nil,err
    }
    pub := pubInterface.(*rsa.PrivateKey)
    signature, err := rsa.SignPKCS1v15(rng, pub, crypto.SHA256, origData[:])
    if err != nil {
        return nil,err
    }
    return signature,nil
}

func (ofd *Orange) ComputeSignature(data string,privatesignKey []byte ) (string,error) {
    buf := bytes.Buffer{}
    binary.Write(&buf, binary.BigEndian, data)
    buf.Write([]byte(data))
    h := sha256.New()
    h.Write(buf.Bytes())
    hash := h.Sum(nil)
    result,err := ofd.RsaEncrypt([]byte(hash),privatesignKey)
    if err != nil {
        return "",err
    }
    str := base64.StdEncoding.EncodeToString(result)
    return str,nil
}

func (ofd *Orange) GenerateDataForCheck(transaction transactionStruct.Transaction)([]map[string]interface{},[]map[string]interface{}){
	var payments  []map[string]interface{}
	var positions  []map[string]interface{}
	entryPayments := make(map[string]interface{})
	entryPositions := make(map[string]interface{})
	for _, product := range transaction.Products {
		quantity := float64(product["quantity"].(float64))
		price := float64(product["value"].(float64))
		entryPayments["type"] = 2
		entryPayments["amount"] = math.Round(quantity*price)
		entryPayments["paymentMethodType"] = 4
		entryPayments["paymentSubjectType"] = 1

		entryPositions["quantity"] = product["quantity"]
		entryPositions["price"] = math.Round(price)
		entryPositions["tax"] = ofd.ConvertTax(product["tax_rate"].(int))
		entryPositions["text"] = product["name"]
		payments = append(payments,entryPayments)
		positions = append(positions,entryPositions)
    }
	return payments,positions
}

func (ofd *Orange) ReadFileCertificate() (string,string,error) {
	crtFilte := fmt.Sprintf("%s%s",pathFileCert,"/ephorOrangeData.crt")
	keyFile := fmt.Sprintf("%s%s",pathFileCert,"/ephorOrangeData.key")
	fcrt, errcrt := os.Open(crtFilte)
    if errcrt != nil {
        return "","",errcrt
    }
    var chunkCrt []byte
    bufcrt := make([]byte, 2048)
	for {
        n, err := fcrt.Read(bufcrt)
        if err != nil && err != io.EOF{
            fmt.Println("read buf fail", err)
            return "","",err
        }
        if n == 0 {
            break
        }
        chunkCrt = append(chunkCrt, bufcrt[:n]...)
    }
	fcrt.Close()
	fkey, errkey := os.Open(keyFile)
    if errkey != nil {
        return "","",errkey
    }
    var chunkkey []byte
    bufkey := make([]byte, 2048)
	for {
        n, err := fkey.Read(bufkey)
        if err != nil && err != io.EOF{
            fmt.Println("read buf fail", err)
            return "","",err
        }
        if n == 0 {
            break
        }
        chunkkey = append(chunkkey, bufkey[:n]...)
    }
	fkey.Close()
	return string(chunkCrt),string(chunkkey),nil
}

func (ofd *Orange) InitData(transaction transactionStruct.Transaction,frModel map[string]interface{})  {
	FrModel = frModel
	ofd.Config.Cert = frModel["auth_public_key"].(string)
	ofd.Config.Key = frModel["auth_private_key"].(string)
	ofd.Config.Group = frModel["param1"].(string)
	ofd.Config.TaxSystem =  transaction.Tax_system
	ofd.Config.Host = frModel["dev_addr"].(string)
	ofd.Config.Port = fmt.Sprintf("%d",frModel["dev_port"])
	ofd.Config.AutomatNumber = transaction.AutomatId
	ofd.Config.Inn = FrModel["inn"].(string)
	TransactionData = transaction
}


func (ofd *Orange) SendCheckApi(data requestFiscal.Data) (map[string]interface{},requestFiscal.Data) {
	ofd.Config.Cert = data.ConfigFR.Cert
	ofd.Config.Key = data.ConfigFR.Key
	ofd.Config.Host = data.ConfigFR.Host
	ofd.Config.Port = data.ConfigFR.Port
	ofd.Config.Inn = data.Inn
	ofd.Config.Sign = data.ConfigFR.Sign
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/", data.ConfigFR.Host, data.ConfigFR.Port)
	Response := ofd.Call("POST", url, data.Fields.Request)
	if Response.Code == 409 {
		result["code"] = Response.Code
		result["fr_id"] = data.CheckId
		result["status"] = "success"
		result["fr_status"] = interfaceFiscal.Status_InQueue
		result["message"] = Response.Message
		return result,data
	}
	if Response.Code != 201 {
		result["code"] = Response.Code
		result["fr_id"] = nil
        result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result,data
	}
	data.DataResponse = Response.Data
	result["code"] = Response.Code
	result["fr_id"] = data.CheckId
	result["status"] = "success"
	result["fr_status"] = interfaceFiscal.Status_InQueue
	result["message"] = "Нет ошибок"
	return result,data
}

func (ofd Orange) GetStatusApi(data requestFiscal.Data) map[string]interface{} {
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/%s/status/%s", ofd.Config.Host, ofd.Config.Port, ofd.Config.Inn, data.CheckId)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	}
	result["code"] = Response.Code
	result["fr_id"] = nil
	result["status"] = "success"
	result["fp"]  = parserTypes.ParseTypeInString(Response.Data["fp"])
	result["fd"] = parserTypes.ParseTypeInFloat64(Response.Data["documentNumber"])
	result["fn"] = parserTypes.ParseTypeInString(Response.Data["fsNumber"])
	result["fr_status"] = interfaceFiscal.Status_Error
	result["message"] = Response.Message
	return result
}

func (ofd *Orange) SendCheck() map[string]interface{} {
	result := make(map[string]interface{})
	key := ""
	TypeFr := int(int64(FrModel["type"].(int64)))
	var dataCheck = make(map[string]interface{})
	var orderString randString.GenerateString
    orderString.RandStringRunes()
	resiptId := orderString.String
	log.Printf("%+v",TransactionData)
	payments,positions := ofd.GenerateDataForCheck(TransactionData)
	content := make(map[string]interface{})
	content["type"] = 1
	content["automatNumber"] = ofd.Config.AutomatNumber
	content["SettlementAddress"] = TransactionData.Address 
	content["SettlementPlace"]	= TransactionData.PointName
	checkClose := make(map[string]interface{})
	checkClose["payments"] = payments
	checkClose["taxationSystem"] = ofd.ConvertTaxationSystem(ofd.Config.TaxSystem)
	content["checkClose"] = checkClose
	content["positions"] = positions
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		key = "4010004"
	}else {
		key = FrModel["inn"].(string)
	}
	dataCheck["id"] = resiptId
	dataCheck["group"] = ofd.Config.Group
	dataCheck["Inn"] = FrModel["inn"]
	dataCheck["key"] = key
	dataCheck["content"] = content
	jsonDataCheck, _ := json.Marshal(dataCheck)
	if TypeFr == interfaceFiscal.Fr_EphorServerOrangeData || TypeFr == interfaceFiscal.Fr_EphorOrangeData {
		certFile,keyFile,errFile := ofd.ReadFileCertificate()
		if errFile != nil {

		}
		ofd.Config.Cert = certFile
		ofd.Config.Key = keyFile
		sign,err := ofd.ComputeSignature(string(jsonDataCheck),privateKeySign);
		if err != nil {

		}
		ofd.Config.Sign = sign
	} else {
		sign,err := ofd.ComputeSignature(string(jsonDataCheck),[]byte(string(FrModel["sign_private_key"].(string))));
		if err != nil {

		}
		ofd.Config.Sign = sign
	}
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/",ofd.Config.Host, ofd.Config.Port)
	Response := ofd.Call("POST", url, jsonDataCheck)
	if Response.Code == 409 {
		return ofd.GetStatus(resiptId)
	}

	if Response.Code != 201 {
		result["code"] = Response.Code
		result["fr_id"] = resiptId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = strings.Join(Response.Errors[:], "\n")
		return result
	}
    return ofd.GetStatus(resiptId)
}

func (ofd *Orange) setTimeOut() (chan bool) {
	timeout := make(chan bool)
	go func() {
			select {
			case <-time.After(5 * time.Minute):
				timeout <- true
			}
	}()
	return timeout
}

func (ofd *Orange) SendRequestOfGetStatus(orderId string) map[string]interface{} {
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s:%s/api/v2/documents/%s/status/%s", ofd.Config.Host, ofd.Config.Port, ofd.Config.Inn, orderId)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code == 200 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "success"
		result["fp"]  = Response.Data["fp"].(string)
		result["fd"] = int(Response.Data["documentNumber"].(float64))
		result["fn"],_ = strconv.Atoi(Response.Data["fsNumber"].(string))
		result["message"] = "нет ошибок"
		result["fr_status"] = interfaceFiscal.Status_Complete
	}
	if Response.Code == 503 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = Response.Message
	}
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = strings.Join(Response.Errors[:], "\n")
	}
	return result
}

func (ofd *Orange) GetStatus(parametrs ...string) map[string]interface{} {
	var orderId string 
	orderId = parametrs[0]
    result := make(map[string]interface{})
	chanTimeOut := ofd.setTimeOut()
    for {
		select {
			case <-time.After(2 * time.Second):
				result = ofd.SendRequestOfGetStatus(orderId)
				if result["status"] == "unsuccess" {
					return result
				}
				if result["code"] == 200 {
					return result
				}
			case <-chanTimeOut:
				result["fr_id"] = orderId
				result["status"] = "unsuccess"
				result["code"] = 0
				result["message"] = fmt.Sprintf("Cancelled by a Timeout of %s", ofd.Name)
				return result
		}
	}
}

func (ofd Orange) Call(method string, url string, json_request []byte) (ResponseOrange) {
	Response := ResponseOrange{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-Signature",ofd.Config.Sign)
	log.Printf("x-Signature %s",ofd.Config.Sign)
	req.Close = true
	log.Println(ofd.Config.Cert)
	log.Println(ofd.Config.Key)
	cert, err := tls.X509KeyPair([]byte(ofd.Config.Cert), []byte(ofd.Config.Key))
	if err != nil {
		Response.Code = 0
		Response.Status = "unsuccess"
		Response.Message = fmt.Sprintf("%v",err)
		return Response
	}
	client := &http.Client{}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	resp, err := client.Do(req)
	if err != nil {
		Response.Code = 0
		Response.Status = "unsuccess"
		Response.Message = fmt.Sprintf("%v",err)
		return Response
	}
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	json.Unmarshal([]byte(body), &Response.Data)
	log.Printf("%+v",Response)
	Response.Code = resp.StatusCode
	if resp.StatusCode > 299 {
		Response.Status = "unsuccess"
		ArrayInterface := Response.Data["errors"].([]interface{})
		errorStrings := parserTypes.ParseArrayInrefaceToArrayString(ArrayInterface)
		Response.Message = strings.Join(errorStrings[:], "\n")
		return Response
	}
	return Response
}

var TransactionData transactionStruct.Transaction
var FrModel map[string]interface{}

func (newf *NewOrangeStruct) NewFiscal() interfaceFiscal.Fiscal  /* тип interfaceFiscal.Fiscal*/ {
    return &NewOrangeStruct{
        Orange: Orange{
        Name: "Orange",
		Config: ConfigOrange{},
       },
    }
}
