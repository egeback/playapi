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
)

// GetJSON downloads contents of url and returns json representation by unmarshal result
func GetJSON(url string) map[string]interface{} {
	req, err := http.NewRequest("GET", url, nil)

	//Fix for SvtPlay
	req.Header.Set("User-Agent", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0) Gecko/20100101 Firefox/77.0")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error getting JSON")
		log.Println(err)
		return make(map[string]interface{})
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return make(map[string]interface{})
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Panic(err)
		resp.Body.Close()
		return result
	}

	return result
}

// Quote string for url use
func Quote(s string) string {
	url := (&url.URL{Path: s}).RequestURI()
	return url
}

//GetStringValue from interface
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

//GetIntValue from interface
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

//Contains checks if string array contains string
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

//ExtractStringSlice gets string slice from interface
func ExtractStringSlice(values []interface{}) []string {
	s := make([]string, 0, len(values))
	for _, value := range values {
		s = append(s, value.(string))
	}
	return s
}

//GetIntValueFromString from interface
func GetIntValueFromString(value string, defaultValue int) *int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return &defaultValue
	}
	return &intValue
}

//GetBoolValueFromString from interface
func GetBoolValueFromString(value string, defaultValue bool) *bool {
	intValue, err := strconv.ParseBool(value)
	if err != nil {
		return &defaultValue
	}
	return &intValue
}
