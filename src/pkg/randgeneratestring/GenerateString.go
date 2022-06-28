package generateString

import(
    "math/rand"
    "time"
    "fmt"
)
type GenerateString struct {
    String string
}
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-")
var orderNumberSize int = 32

func (gen *GenerateString) getDataOfTimeString() string {
    today := time.Now()
    hour := today.Hour()
	minute := today.Minute()
	second := today.Second()
    TimeString := fmt.Sprintf("%d%d%d", hour, minute, second)
    return TimeString
}


func (gen *GenerateString)RandStringRunes() {
    stringResult := ""
    b := make([]rune, orderNumberSize)
    for k := 0; k < len(b); k++ {
        if k == 0 && k < 3 {
            stringResult += string(letterRunes[rand.Intn(len(letterRunes))])
        }
    }
    stringResult += "-"
    Time := gen.getDataOfTimeString()
    stringResult += Time
    stringResult += "-"
    for k := 0; k < len(b); k++ {
        if k == 0 && k < 3 {
            stringResult += string(letterRunes[rand.Intn(len(letterRunes))])
        }
    }  
//     for i := range b {
//       b[i] = letterRunes[rand.Intn(len(letterRunes))] 
//     }
//    gen.String = string(b)
    gen.String = stringResult

}
