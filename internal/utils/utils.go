package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// GetJSON ...
func GetJSON(url string) map[string]interface{} {
	resp, err := http.Get(url)

	if err != nil {
		log.Panic(err)
		return make(map[string]interface{})
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
		return make(map[string]interface{})
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Panic(err)
		resp.Body.Close()
		return result
	}
	resp.Body.Close()

	return result
}

// Quote ...
func Quote(s string) string {
	url := (&url.URL{Path: s}).RequestURI()
	return url
	//strings.ReplaceAll(url, ",", "%2C")
}

type tagOptions string

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// ParseTag ...
func ParseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// IsEmptyValue ...
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func GetStringValue(data map[string]interface{}, param string, def string) string {
	value := def
	if _, ok := data[param]; ok {
		switch data[param].(type) {
		case string:
			value = data[param].(string)
		case int:
			value = string(data[param].(int))
		case float64:
			if data[param].(float64) == math.Trunc(data[param].(float64)) {
				value = fmt.Sprintf("%.0f", data[param].(float64))
			} else {
				value = fmt.Sprintf("%f", data[param].(float64))
			}
		case nil:

		default:
			fmt.Println("type:", reflect.TypeOf(data[param]))
		}

	}
	return value
}

func GetIntValue(data map[string]interface{}, param string, def int) int {
	value := def
	if _, ok := data[param]; ok {
		switch data[param].(type) {
		case string:
			val, err := strconv.Atoi(data[param].(string))
			if err != nil {
				log.Panic(err)
				return def
			}
			value = val
		case int:
			value = data[param].(int)
		}

	}
	return value
}
