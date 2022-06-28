package fermaOfd

import (
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
    "bytes"
	"encoding/json"
	"encoding/base64"
	"math"
	"time"
	randString "ephorservices/src/pkg/randgeneratestring"
    transactionStruct "ephorservices/src/data/transaction"
    interfaceFiscal "ephorservices/src/server/utils/interface/fiscal"
	requestFiscal "ephorservices/src/data/requestApi"
	parserTypes "ephorservices/src/pkg/parser/typeParse"
	"regexp"
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

type ResponseOfd struct {
	Code   int
	Status string
	Message string
	Error  struct {
		Code    int
		Message string
	}
	Data   map[string]interface{}
}

func (r ResponseOfd) GetDataString(field string) string {
	str, _ := r.Data[field].(string)
	return str
}

func (r ResponseOfd) GetDataMap(field string) map[string]interface{} {
	i, _ := r.Data[field].(map[string]interface{})
	return i
}

type ConfigOfd struct {
	Auth struct {
		Token string
	}
	ReceiptId string
	Host, Login, Password,Inn string
	TaxSystem,AutomatNumber int
}

type FermaOfd struct {
    Name string
	Config ConfigOfd
}

type NewFermaOfdStruct struct {
    FermaOfd
}

func (ofd *FermaOfd) ConvertTax(tax int) string {
	switch tax {
		 case TaxRate_NDSNone:
		 return "VatNo"
		 fallthrough
		 case TaxRate_NDS0:
		 return "Vat0"
		 fallthrough
		 case TaxRate_NDS10:
		 return "Vat10"
		 fallthrough
		 case TaxRate_NDS18:
		 return "Vat20"
	 }
	 return "VatNo"
}

func (ofd *FermaOfd) ConvertTaxationSystem(taxsystem int) string {
	switch taxsystem {
		 case TaxSystem_OSN:
		 return "Common"
		 fallthrough
		 case TaxSystem_USND:
		 return "SimpleIn"
		 fallthrough
		 case TaxSystem_USNDMR:
		 return "SimpleInOut"
		 fallthrough
		 case TaxSystem_ENVD:
		 return "Unified"
		 fallthrough
		 case TaxSystem_ESN:
		 return "UnifiedAgricultural"
		 fallthrough
		 case TaxSystem_Patent:
		 return "Patent"
	 }
	 return "Common"
}

func (ofd *FermaOfd) MakeUrlQr(date string, summ int, frResponse map[string]interface{}) string {
	t, _ := time.Parse(layoutISO, date)
	valueSumm := summ/100
    stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1",fmt.Sprintf("%s",t.Format(layoutQr)),fmt.Sprintf("%v.00",valueSumm),frResponse["fn"],frResponse["fd"],frResponse["fp"]);
	log.Println(stringResult)
	return stringResult
}

func (ofd *FermaOfd) GetQrPicture(date string, summ int, frResponse map[string]interface{}) string {
	result := ofd.MakeUrlQr(date,summ,frResponse)
	str := base64.StdEncoding.EncodeToString([]byte(result))
	return str
}

func (ofd *FermaOfd) GetQrUrl(date string, summ int, frResponse map[string]interface{}) string {
	return ofd.MakeUrlQr(date,summ,frResponse)
}

func (ofd *FermaOfd) ConvertPaymentDevice(paymentDevice string) int {
	switch paymentDevice {
		case "CA":
			return 0
			fallthrough
		case "DA":
			return 1
	}
	return 0
}

func (ofd *FermaOfd) MakeAuth() map[string]interface{} {
	result := make(map[string]interface{})
	url := fmt.Sprintf("https://%s/api/Authorization/CreateAuthToken", ofd.Config.Host)
	json_str := fmt.Sprintf(`{"Login":"%s", "Password": "%s"}`, ofd.Config.Login, ofd.Config.Password)
	Response := ofd.Call("POST", url, []byte(json_str))
	if Response.Code != 200 {
		result["status"] = "unsuccess"
		return result
	}
	ofd.Config.Auth.Token = Response.GetDataString("AuthToken")
	if len(ofd.Config.Auth.Token) == 0 {
		result["status"] = "unsuccess"
		result["message"] = Response.Error.Message
		return result
	}
	result["status"] = "success"
	return result
}

func (ofd *FermaOfd) GenerateDataForCheck(transaction transactionStruct.Transaction)([]map[string]interface{},[]map[string]interface{}){
	var summ int
	var payments  []map[string]interface{}
	var positions  []map[string]interface{}
	entryPayments := make(map[string]interface{})
	entryPositions := make(map[string]interface{})
	entryPayments["PaymentType"] = 1
	entryPayments["Sum"] = 0
	for _, product := range transaction.Products {
		quantity := float64(product["quantity"].(float64))
		price := float64(product["value"].(float64))
		summ += int(math.Round(quantity*price))
		entryPositions["PaymentMethod"] = 4
		entryPositions["PaymentType"] = 1
		entryPositions["Quantity"] = product["quantity"]
		entryPositions["Price"] = math.Round(price)
		entryPositions["Amount"] = math.Round(quantity*price)
		entryPositions["Vat"] = ofd.ConvertTax(product["tax_rate"].(int))
		entryPositions["Label"] = product["name"]
		positions = append(positions,entryPositions)
    }
	entryPayments["Sum"] = summ
	payments = append(payments,entryPayments)
	return payments,positions
}

func (ofd *FermaOfd) InitData(transaction transactionStruct.Transaction,frModel map[string]interface{})  {
	FrModel = frModel
	ofd.Config.TaxSystem =  transaction.Tax_system
	ofd.Config.Host = frModel["dev_addr"].(string)
	ofd.Config.AutomatNumber = transaction.AutomatId
	ofd.Config.Inn = FrModel["inn"].(string)
	TransactionData = transaction
}


func (ofd *FermaOfd) SendCheckApi(data requestFiscal.Data) (map[string]interface{},requestFiscal.Data) {
	result := make(map[string]interface{})
	ofd.Config.Login = data.ConfigFR.Login
	ofd.Config.Password = data.ConfigFR.Password
	ofd.Config.Host = data.ConfigFR.Host
	ofd.Config.Inn = data.Inn
	resultAuth := ofd.MakeAuth()
	if resultAuth["status"] == "unsuccess" {
			result["code"] = 400
			result["fr_id"] = nil
			result["fp_string"] = nil
			result["status"] = "unsuccess"
			result["fr_status"] = interfaceFiscal.Status_Error
			result["message"] = ""
			return result,data
	}
	url := fmt.Sprintf("https://%s/api/kkt/cloud/receipt?AuthToken=%s", ofd.Config.Host, ofd.Config.Auth.Token)
	Response := ofd.Call("POST", url, data.Fields.Request)
	if Response.Code != 200 {
		if Response.Error.Code == 1019 {
			re := regexp.MustCompile("--([a-f0-9-]+)--")
			match := re.FindStringSubmatch(Response.Error.Message)
			if len(match) > 1 {
				result["code"] = 200
				result["fr_id"] = match[1]
				result["status"] = "success"
				result["fr_status"] = interfaceFiscal.Status_InQueue
				result["message"] = "Нет ошибок"
				data.DataResponse["ReceiptId"] = match[1]
				return result,data
			}
			result["code"] = Response.Code
			result["fr_id"] = nil
			result["fp_string"] = nil
			result["status"] = "unsuccess"
			result["fr_status"] = interfaceFiscal.Status_Error
			result["message"] = "Ошибка"
			return result,data
		}
	}
	data.DataResponse["ReceiptId"] = Response.GetDataString("ReceiptId")
	if len(data.DataResponse["ReceiptId"].(string)) == 0 {
		result["code"] = 400
		result["fr_id"] = nil
		result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = fmt.Sprintf("Error: %s No ReceiptId", ofd.Name)
		return result,data
	}
	result["code"] = Response.Code
	result["fr_id"] = data.DataResponse["ReceiptId"].(string)
	result["status"] = "success"
	result["fr_status"] = interfaceFiscal.Status_InQueue
	result["message"] = "Нет ошибок"
	return result,data
}

func (ofd FermaOfd) GetStatusApi(data requestFiscal.Data) map[string]interface{} {
	result := make(map[string]interface{})
	json_str := fmt.Sprintf(`{"Request":{"ReceiptId": "%s"}}`, data.DataResponse["ReceiptId"])
	url := fmt.Sprintf("https://%s/api/kkt/cloud/status?AuthToken=%s", ofd.Config.Host, ofd.Config.Auth.Token)
	Response := ofd.Call("POST", url, []byte(json_str))
	status := Response.GetDataString("StatusName")
	if status == "CONFIRMED" {
		device := Response.GetDataMap("Device")
		result["code"] = Response.Code
		result["fr_id"] = data.DataResponse["ReceiptId"]
		result["status"] = "success"
		result["fp"]  = parserTypes.ParseTypeInString(device["FPD"])
		result["fd"] = parserTypes.ParseTypeInString(device["FDN"])
		result["fn"] = parserTypes.ParseTypeInString(device["FN"])
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	} else if status == "KKT_ERROR" {
		result["code"] = Response.Code
		result["fr_id"] = data.DataResponse["ReceiptId"]
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	}
	result["status"] = "success"
    result["code"] = 202
    return result
}

func (ofd *FermaOfd) SendCheck() map[string]interface{} {
	result := make(map[string]interface{})
	var dataCheck = make(map[string]interface{})
	var orderString randString.GenerateString
    orderString.RandStringRunes()
	resiptGenId := orderString.String
	date := TransactionData.Date
	if date == "" {
		t := time.Now() 
		date = t.Format("2006-01-02 15:04:05")
	}
	payments,positions := ofd.GenerateDataForCheck(TransactionData)
	content := make(map[string]interface{})
	content["Inn"] = FrModel["inn"]
	content["Type"] = "Income"
	content["InvoiceId"] = resiptGenId
	content["LocalDate"] = date
	content["SettlementAddress"] = TransactionData.Address 
	content["SettlementPlace"]	= TransactionData.PointName
	CustomerReceipt := make(map[string]interface{})
	CustomerReceipt["TaxationSystem"] = ofd.ConvertTaxationSystem(ofd.Config.TaxSystem)
	CustomerReceipt["Email"] = FrModel["email"]
	CustomerReceipt["Phone"] = FrModel["phone"]
	CustomerReceipt["PaymentType"] = 1
	CustomerReceipt["AutomatNumber"] = ofd.Config.AutomatNumber
	CustomerReceipt["KktFA"] = true
	CustomerReceipt["BillAddress"] = TransactionData.Address
	CustomerReceipt["Items"] = positions
	CustomerReceipt["PaymentItems"] = payments
	content["CustomerReceipt"] = CustomerReceipt
	dataCheck["Request"] = content
	jsonDataCheck, _ := json.Marshal(dataCheck)
	url := fmt.Sprintf("https://%s/api/kkt/cloud/receipt?AuthToken=%s", ofd.Config.Host, ofd.Config.Auth.Token)
	Response := ofd.Call("POST", url, jsonDataCheck)
	if Response.Code != 200 {
		if Response.Error.Code == 1019 {
			re := regexp.MustCompile("--([a-f0-9-]+)--")
			match := re.FindStringSubmatch(Response.Error.Message)
			if len(match) > 1 {
				result["code"] = 200
				result["fr_id"] = match[1]
				result["status"] = "success"
				result["fr_status"] = interfaceFiscal.Status_InQueue
				result["message"] = "Нет ошибок"
				return result
			}
			result["code"] = Response.Code
			result["fr_id"] = nil
			result["fp_string"] = nil
			result["status"] = "unsuccess"
			result["fr_status"] = interfaceFiscal.Status_Error
			result["message"] = "Ошибка"
			return result
		}
	}
	resiptId := Response.GetDataString("ReceiptId")
	if len(Response.GetDataString("ReceiptId")) == 0 {
		result["code"] = 400
		result["fr_id"] = nil
		result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = fmt.Sprintf("Error: %s No ReceiptId", ofd.Name)
		return result
	}
    return ofd.GetStatus(resiptId)
}

func (ofd *FermaOfd) setTimeOut() (chan bool) {
	timeout := make(chan bool)
	go func() {
			select {
			case <-time.After(5 * time.Minute):
				timeout <- true
			}
	}()
	return timeout
}

func (ofd *FermaOfd) SendRequestOfGetStatus(orderId string) map[string]interface{} {
	result := make(map[string]interface{})
	json_str := fmt.Sprintf(`{"Request":{"ReceiptId": "%s"}}`, orderId)
	url := fmt.Sprintf("https://%s/api/kkt/cloud/status?AuthToken=%s", ofd.Config.Host, ofd.Config.Auth.Token)
	Response := ofd.Call("POST", url, []byte(json_str))
	status := Response.GetDataString("StatusName")
	if status == "CONFIRMED" {
		device := Response.GetDataMap("Device")
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "success"
		result["fp"]  = parserTypes.ParseTypeInString(device["FPD"])
		result["fd"] = parserTypes.ParseTypeInString(device["FDN"])
		result["fn"] = parserTypes.ParseTypeInString(device["FN"])
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	} else if status == "KKT_ERROR" {
		result["code"] = Response.Code
		result["fr_id"] = orderId
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Message
		return result
	}
    return result
}

func (ofd *FermaOfd) GetStatus(parametrs ...string) map[string]interface{} {
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

func (ofd FermaOfd) Call(method string, url string, json_request []byte) (ResponseOfd) {
	Response := ResponseOfd{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	client := &http.Client{}
	
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
	Response.Code = resp.StatusCode
	if resp.StatusCode > 299 {
		Response.Status = "unsuccess"
		Response.Message = Response.Error.Message
		return Response
	}
	return Response
}

var TransactionData transactionStruct.Transaction
var FrModel map[string]interface{}

func (newf *NewFermaOfdStruct) NewFiscal() interfaceFiscal.Fiscal  /* тип interfaceFiscal.Fiscal*/ {
    return &NewFermaOfdStruct{
        FermaOfd: FermaOfd{
        Name: "FermaOfd",
		Config: ConfigOfd{},
       },
    }
}
