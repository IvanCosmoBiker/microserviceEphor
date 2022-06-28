package typeParse

import (
    "strconv"
    "fmt"
)

func ParseTypeInFloat64(parametr interface{}) float64 {
    switch parametr.(type) {
        case string:
           value,_ := strconv.ParseFloat(parametr.(string), 64)
           return value
        case int:
            return float64(parametr.(int))
        case int8:
            return float64(parametr.(int8)) 
        case int16:
            return float64(parametr.(int16))
        case int32:
            return float64(parametr.(int32)) 
        case int64:
            return float64(parametr.(int64))
        case uint8:
            return float64(parametr.(uint8)) 
        case uint16:
            return float64(parametr.(uint16))
        case uint32:
            return float64(parametr.(uint32))
        case uint64:
            return float64(parametr.(uint64))
        case float32:
            return float64(parametr.(float32))
    }
    return parametr.(float64)
}

func ParseTypeInString(parametr interface{}) string {
    switch parametr.(type) {
        case int,int8,int16,int32,int64,uint8,uint16,uint32,uint64,complex64,complex128,float32,float64:
            return fmt.Sprintf("%v",parametr)
    }
    return parametr.(string)
}

func ParseArrayInrefaceToArrayString(parametr []interface{}) []string {
    s := make([]string, len(parametr))
    for i, v := range parametr {
        s[i] = fmt.Sprint(v)
    }
    return s
}

