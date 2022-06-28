package nanokass

import (
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
    "bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/rand"
	"encoding/json"
	"encoding/pem"
	"encoding/base64"
	"math"
	"errors"
	"strings"
	"time"
    "crypto/des"
    "encoding/hex"
    "crypto/cipher"
    "crypto/hmac"
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

var RSA_PUB_FIRST = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwFXHnzc5YKj8e3tlNzST
CkA8Tq4gjTH0VMuhJhg5QWpFjFKwtnK3u4EOaQGmjqDtzyffVHmKuGikg9jE20sG
nJN4hTtySihOiUWRd4zhJVMevBQmsEQS33bg26UzzKCeO12mbM/Q4ip7YXEfWM/F
Tq2l94psQgmIDh/LtHVf3OBlz8I6u5VaP3AS0Hv9RBUin0RBkRUC+5tgURm382XT
nJ2GzZ8cEGJm3C+s0+W1N2igjV0X3MihylHGDyl+8FpbFIlXsaJOYQ0//JIgnaBz
MV2JyNTHBzPJrcIMHIbKBVAmDLfgeDNKug7wIadEcqoJaCz74yG9l9nJWISWQkI6
Ed8nDVsoaIkMQBuWWxfHjQEU8R8OVjRzhOGHPG2ka6y1/jcOS5JWPzS5YVXRPbrh
QYcoNebsOBaFxJYZ2E7VhVdrGWlBqhANFba7umZXVOvmDXIsH974Yv4awAaP70VP
SLFIdjiNy/SB8w0O8PJOUPznpMhvi1clBgp3PvtYmhUqmdHWPwjcjy0JmY9KrWz0
0Im1yDTTybtV3uYnwR677TmsLmR9c6T7EHlT3gG6Y0bM3w9tyrGqVKy1jIkyUZPV
f0dmXTfbh+hcC5kYal+M7lcn7wSSLHTUk+C/YWE1e5TvTBK6teU2VNmz80Yt2IS2
mcXlfKlZXilMmPJCdUI7nNMCAwEAAQ==
-----END PUBLIC KEY-----`)

var RSA_PUB_SECOND = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA+fu+NGlnWAXqIVgEL37v
eatlyooYi+iHLiBmCDowNZUBAiQ+pvbnzkowUKdr86lGrzQLCAvVyXWG0U4kdixA
X0GTkIR/3g3h2/8hRx0x3K0umT+tcZC3iJytKzP+EM/B6sDdw6/URbykwvrAlbQs
G9d6eCqq0F/6muOM3gQazy8CuHyx4iFQpml4E1/IQgp3tZJOX5I9xieHTUwct2Ok
URCKYnHJZrRIN9rwXQkNG1q+M8HDqI1Mwq88wieVC+SUuoPc8F0MlIWs2zwDhLcX
84OQTRFqlW3NFR/6kUn3TIC1JZD1Ft/8fWukZzAFsAmdXmFzhBUuBPvjIzzLafY3
f8IszADMnloJ0BW3iGVRGj6hygX7Jpr/86LPHu6PBJzHzCp9bnfOiSjRENzzy55f
DdVbYpVgWDt4+UEkl9qNRNuiSMDpKeVNy6jxbihZneYCR8alnH8Olh6lL7bmGdww
qI9LSyq/qFfIMDV8onit/dLxzypFJofRfjZ1Dc8ZEqh2sab8qEMNPGQwTM/FVFWM
bq0hmjjY+BFWGY/h0z1NZMX75Uzyd9OdXaRoTlHPfOxxAIfclP2XY2K8f5PQ37g/
fX2R8bw/fXQd2ndi/+uPCGK92Xw4/3/osJKpm3QSYhSda53T9Ddned7BtWDQJqdV
Y/SUskwLLyjtSb0LqsSKBHkCAwEAAQ==
-----END PUBLIC KEY-----`)
      
    const HMAC_FIRST  = "BBuXaXBdHg+wLPjRJpf3N/NmLq5kuvzGQx3II15/j8o="
    const HMAC_SECOND = "aFZP3PbvrMZNNxxqJxaCnCLama5L8H1/YGO3UYsoCVQ="
    const URL_TO_SEND_TO_NANOKASSA = "http://q.nanokassa.ru/srv/igd.php"

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

type ResponseNanokass struct {
	Code   int
	Status string
	Message string
    Data map[string]interface{}
}

type ConfigNanokass struct {
	Cert string
	Host string
	Port string 
	Key string
	Group string 
	TaxSystem int
	Sign string
	AutomatNumber int
	Inn string
    Sign_private_key string
}

type Nanokass struct {
    Name string
}

type NewNanokassStruct struct {
    Nanokass
}

func (ofd *Nanokass) ConvertTax(tax int) int {
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

func (ofd *Nanokass) ConvertTaxationSystem(taxsystem int) int {
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

func (ofd *Nanokass) MakeUrlQr(date string, summ int, frResponse map[string]interface{}) string {
	t, _ := time.Parse(layoutISO, date)
	valueSumm := summ/100
    stringResult := fmt.Sprintf("t=%s&s=%v&fn=%v&i=%v&fp=%v&n=1",fmt.Sprintf("%s",t.Format(layoutQr)),fmt.Sprintf("%v.00",valueSumm),frResponse["fn"],frResponse["fd"],frResponse["fp"]);
	log.Println(stringResult)
	return stringResult
}

func (ofd *Nanokass) GetQrPicture(date string, summ int, frResponse map[string]interface{}) string {
	result := ofd.MakeUrlQr(date,summ,frResponse)
	str := base64.StdEncoding.EncodeToString([]byte(result))
	return str
}

func (ofd *Nanokass) GetQrUrl(date string, summ int, frResponse map[string]interface{}) string {
	return ofd.MakeUrlQr(date,summ,frResponse)
}

func (ofd *Nanokass) GenerateRandomBytesInString(n int) string {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return ""
    }
    return base64.URLEncoding.EncodeToString(b)
}

func DesCBCEncrypt(data, key, iv []byte) ([]byte, error) {
    block, err := des.NewCipher(key)
    if err != nil {
        return nil, err
    }

    data = Pkcs5Padding(data, block.BlockSize())
    cryptText := make([]byte, len(data))

    blockMode := cipher.NewCBCEncrypter(block, iv)
    blockMode.CryptBlocks(cryptText, data)
    return cryptText, nil
}

func Pkcs5Padding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padText...)
}

func (ofd *Nanokass) RsaEncrypt(origData,publickey []byte) ([]byte, error) {
    rng := rand.Reader
    block, _ := pem.Decode(publickey)
    if block == nil {
        return nil, errors.New("is not pem")
    }
    pubInterface,err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        return nil,err
    }
    pub := pubInterface
    signature, err := rsa.EncryptPKCS1v15(rng, pub, origData[:])
    if err != nil {
        return nil,err
    }
    return signature,nil
}

func (ofd *Nanokass) Crypt_nanokassa_first(data []byte) map[string]interface{} {
    resultMap := make(map[string]interface{})
    IVdata := ofd.GenerateRandomBytesInString(16)
    pw := ofd.GenerateRandomBytesInString(32)
    mk := HMAC_FIRST
    result,err := DesCBCEncrypt(data,[]byte(pw),[]byte(IVdata))
    if err != nil {
        fmt.Println(err)
        return resultMap
    }
    
    h := hmac.New(sha256.New, []byte(mk))
    h.Write([]byte(IVdata + string(result)))
    hexData := hex.EncodeToString(h.Sum(nil))
    returnDataDE := base64.StdEncoding.EncodeToString([]byte(IVdata + string(result) + hexData))
    ab_rsa,_ := ofd.RsaEncrypt([]byte(pw),RSA_PUB_FIRST)
    returnDataAB := base64.StdEncoding.EncodeToString(ab_rsa)
    resultMap["ab"] = returnDataAB
    resultMap["de"] = returnDataDE
    return resultMap
}

func (ofd *Nanokass) Crypt_nanokassa_second(data []byte) map[string]interface{} {
    resultMap := make(map[string]interface{})
    IVdata := ofd.GenerateRandomBytesInString(16)
    pw := ofd.GenerateRandomBytesInString(32)
    mk := HMAC_SECOND
    result,err := DesCBCEncrypt(data,[]byte(pw),[]byte(IVdata))
    if err != nil {
        fmt.Println(err)
        return resultMap
    }
    
    h := hmac.New(sha256.New, []byte(mk))
    h.Write([]byte(IVdata + string(result)))
    hexData := hex.EncodeToString(h.Sum(nil))
    returnDataDEE := base64.StdEncoding.EncodeToString([]byte(IVdata + string(result) + hexData))
    aab_rsa,_ := ofd.RsaEncrypt([]byte(pw),RSA_PUB_SECOND)
    returnDataAAB := base64.StdEncoding.EncodeToString(aab_rsa)
    resultMap["aab"] = returnDataAAB
    resultMap["dde"] = returnDataDEE
    return resultMap
}

func (ofd *Nanokass) GenerateDataForCheck(transaction transactionStruct.Transaction)([]map[string]interface{},[]map[string]interface{},int){
    var summ int
	var payments  []map[string]interface{}
	var positions  []map[string]interface{}
	entryPayments := make(map[string]interface{})
	entryPositions := make(map[string]interface{})
	entryPayments["dop_rekvizit_1192"] = ""
	entryPayments["inn_pokupatel"] = ""
	entryPayments["name_pokupatel"] = ""
	entryPayments["rezhim_nalog"] = ofd.ConvertTaxationSystem(transaction.Tax_system)
	entryPayments["kassir_inn"] = ""
	entryPayments["kassir_fio"] = ""
	entryPayments["client_email"] = "none"
	entryPayments["money_nal"] = 0
	entryPayments["money_predoplata"] = 0
	entryPayments["money_postoplata"] = 0
	entryPayments["money_vstrecha"] = 0
	entryPayments["money_electro"] = 0
	for _, product := range transaction.Products {
		quantity := float64(product["quantity"].(float64))
		price := float64(product["value"].(float64))
        summ += int(math.Round(quantity*price))
		entryPayments["money_electro"] = int(math.Round(quantity*price))
        entryPositions["summa"] = math.Round(quantity*price)
        entryPositions["price_piece_bez_skidki"] = math.Round(price)
        entryPositions["priznak_sposoba_rascheta"] = 4
        entryPositions["priznak_predmeta_rascheta"] = 1
		entryPositions["kolvo"] = product["quantity"]
        entryPositions["name_tovar"] = product["name"]
		entryPositions["price_piece"] = math.Round(price)
		entryPositions["stavka_nds"] = ofd.ConvertTax(product["tax_rate"].(int))
		entryPositions["priznak_agenta"] = "none"
		positions = append(positions,entryPositions)
    }
	entryPayments["money_electro"] = summ
	payments = append(payments,entryPayments)
	return payments,positions,summ
}


var TransactionData transactionStruct.Transaction
var FrModel map[string]interface{}
var Config ConfigNanokass 

func (ofd *Nanokass) InitData(transaction transactionStruct.Transaction,frModel map[string]interface{})  {
	Config = ConfigNanokass{}
	FrModel = frModel
	Config.Cert = frModel["auth_public_key"].(string)
	Config.Key = frModel["auth_private_key"].(string)
	Config.Group = frModel["param1"].(string)
	Config.TaxSystem =  transaction.Tax_system
	Config.Host = frModel["dev_addr"].(string)
	Config.Port = fmt.Sprintf("%d",frModel["dev_port"])
	Config.AutomatNumber = transaction.AutomatId
	Config.Inn = FrModel["inn"].(string)
    Config.Sign_private_key = FrModel["sign_private_key"].(string)
	TransactionData = transaction
}

func (ofd *Nanokass) SendCheckApi(data requestFiscal.Data) (map[string]interface{},requestFiscal.Data) {
	result := make(map[string]interface{})
	url := URL_TO_SEND_TO_NANOKASSA
	Response := ofd.Call("POST", url, data.Fields.Request)
	if Response.Code != 200 {
		result["code"] = Response.Code
		result["fr_id"] = nil
        result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = "Ошибка"
		return result,data
	}
	if Response.Data["status"] != "success" {
		result["code"] = Response.Code
		result["fr_id"] = nil
        result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = Response.Data["error"]
		return result,data
	}
	data.DataResponse = Response.Data
	result["code"] = Response.Code
	result["fr_id"] = Response.Data["qnuid"]
    result["fp_string"] = Response.Data["nuid"]
	result["status"] = "success"
	result["fr_status"] = interfaceFiscal.Status_InQueue
	result["message"] = "Нет ошибок"
    return result,data
}

func (ofd *Nanokass) GetStatusApi(data requestFiscal.Data) map[string]interface{} {
	result := make(map[string]interface{})
	nuid := data.DataResponse["nuid"].(string)
    qnuid := data.DataResponse["qnuid"].(string)
	url := fmt.Sprintf("http://fp.nanokassa.com/getfp?nuid=%s&qnuid=%s&auth=base", nuid, qnuid)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = qnuid
        result["fp_string"] = nuid
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = "Ошибка"
	}
	if int(Response.Data["check_status"].(float64)) == 0 {
            result["status"] = "success"
            result["code"] = 202
            return result
    }
	if int(Response.Data["check_status"].(float64)) == 1 ||  int(Response.Data["check_status"].(float64)) == 3 {
		result["code"] = Response.Code
		result["fr_id"] = qnuid
		result["fp_string"] = nuid
		result["status"] = "success"
		result["fp"]  = parserTypes.ParseTypeInString(Response.Data["check_num_fp"])
		result["fd"] = parserTypes.ParseTypeInFloat64(Response.Data["check_num_fd"])
		result["fn"] = parserTypes.ParseTypeInString(Response.Data["check_fn_num"])
		result["message"] = "нет ошибок"
		result["fr_status"] = interfaceFiscal.Status_Complete
		return result
	}
	return result

}

func (ofd *Nanokass) SendCheck() map[string]interface{} {
	result := make(map[string]interface{})
	log.Printf("%+v",TransactionData)
	payments,positions,summ := ofd.GenerateDataForCheck(TransactionData)
    paramsNano := Config.Sign_private_key
    dataNano := strings.Split(paramsNano, ":")
	content1 := make(map[string]interface{})
	content1["kassaid"] = dataNano[0]
	content1["kassatoken"] = dataNano[1]
	content1["cms"] = "wordpress"  
	content1["check_send_type"]	= "email" 
    content1["check_vend_address"] = TransactionData.Address
    content1["check_vend_mesto"] = TransactionData.PointName
    content1["check_vend_num_avtovat"] = TransactionData.AutomatId
    content1["products_arr"] = positions
    content1["oplata_arr"] = payments
    itog := make(map[string]interface{})
    itog["itog_cheka"] = summ
    itog["priznak_rascheta"] = 1
    content1["itog_arr"] = itog
    request1,_ := json.Marshal(content1)
    firstcrypt := ofd.Crypt_nanokassa_first(request1)
    returnDataAB := firstcrypt["ad"]
	returnDataDE := firstcrypt["de"]
    content2 := make(map[string]interface{})
    content2["ab"] = fmt.Sprintf("'%v'",returnDataAB)
    content2["de"] = fmt.Sprintf("'%v'",returnDataDE)
    content2["kassaid"] = fmt.Sprintf("'%s'",dataNano[0])
    content2["kassatoken"] = fmt.Sprintf("'%s'",dataNano[1])
    content2["check_type"] = "standart"
    content2["test"] = "0"
    request2,_ := json.Marshal(content2)
    secondcrypt := ofd.Crypt_nanokassa_second(request2)
    returnDataAAB := secondcrypt["aab"]
	returnDataDE2 := secondcrypt["dde"]
    content3 := make(map[string]interface{})
    content3["aab"] = fmt.Sprintf("'%v'",returnDataAAB)
    content3["aab"] = fmt.Sprintf("'%v'",returnDataDE2)
    content3["test"] = "0"
    jsonDataCheck,_ :=  json.Marshal(content3)
	url := URL_TO_SEND_TO_NANOKASSA
	Response := ofd.Call("POST", url, jsonDataCheck)

	if Response.Code != 200 {
		result["code"] = Response.Code
		result["fr_id"] = nil
		result["fp_string"] = nil
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = ""
		return result
	}
    nuid := Response.Data["nuid"].(string)
    qnuid := Response.Data["qnuid"].(string)
    return ofd.GetStatus(nuid,qnuid)
}

func (ofd *Nanokass) setTimeOut() (chan bool) {
	timeout := make(chan bool)
	go func() {
			select {
			case <-time.After(5 * time.Minute):
				timeout <- true
			}
	}()
	return timeout
}

func (ofd *Nanokass) SendRequestOfGetStatus(nuid,qnuid string) map[string]interface{} {
	result := make(map[string]interface{})
    
	url := fmt.Sprintf("http://fp.nanokassa.com/getfp?nuid=%s&qnuid=%s&auth=base", nuid, qnuid)
	Response := ofd.Call("GET", url, []byte(""))
	if Response.Code == 200 {
        if Response.Data["check_status"] == "0" {
            result["status"] = "success"
            result["code"] = 202
            return result
        }
        if Response.Data["check_status"] == "1" ||  Response.Data["check_status"] == "3" {
            result["code"] = Response.Code
            result["fr_id"] = qnuid
            result["fp_string"] = nuid
            result["status"] = "success"
            result["fp"]  = Response.Data["check_num_fp"].(string)
            result["fd"] = Response.Data["check_num_fd"].(string)
            result["fn"],_ = Response.Data["check_fn_num"].(string)
            result["message"] = "нет ошибок"
            result["fr_status"] = interfaceFiscal.Status_Complete
        }else {
            result["status"] = "success"
            result["code"] = 202
            return result
        }
	}
	if Response.Code == 503 {
		result["code"] = Response.Code
		result["fr_id"] = qnuid
        result["fp_string"] = nuid
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Overflow
		result["message"] = ""
	}
	if Response.Code > 299 {
		result["code"] = Response.Code
		result["fr_id"] = qnuid
        result["fp_string"] = nuid
		result["status"] = "unsuccess"
		result["fr_status"] = interfaceFiscal.Status_Error
		result["message"] = "Ошибка"
	}
	return result
}

func (ofd *Nanokass) GetStatus(parametrs ...string) map[string]interface{} {
    var nuid string
    var qnuid string
    nuid = parametrs[0]
    qnuid = parametrs[1]
    result := make(map[string]interface{})
	chanTimeOut := ofd.setTimeOut()
    for {
		select {
			case <-time.After(2 * time.Second):
				result = ofd.SendRequestOfGetStatus(nuid,qnuid)
				if result["status"] == "unsuccess" {
					return result
				}
				if result["code"] == 200 {
					return result
				}
			case <-chanTimeOut:
				result["fr_id"] = qnuid
                result["fp_string"] = nuid
				result["status"] = "unsuccess"
				result["code"] = 0
				result["message"] = fmt.Sprintf("Cancelled by a Timeout of %s", ofd.Name)
				return result
		}
	}
}

func (ofd *Nanokass) Call(method string, url string, json_request []byte) (ResponseNanokass) {
	Response := ResponseNanokass{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(json_request))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
	log.Printf("%+v",Response)
	Response.Code = resp.StatusCode
	if resp.StatusCode > 299 {
		Response.Status = "unsuccess"
		Response.Message = "ошибка"
		return Response
	}
	return Response
}


func (newf *NewNanokassStruct) NewFiscal() interfaceFiscal.Fiscal  /* тип interfaceFiscal.Fiscal*/ {
    return &NewNanokassStruct{
        Nanokass: Nanokass{
        Name: "Nanokass",
       },
    }
}
