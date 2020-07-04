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
	"time"
)

//GetTimeFromString from string
func GetTimeFromString(str string, layouts ...string) *time.Time {
	if len(layouts) == 0 {
		layouts = []string{"2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05Z"}
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, str)
		if err != nil {
			continue
		}
		return &t
	}
	return nil
}

//Min of two ints
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// GetJSON downloads contents of url and returns json representation by unmarshal result
func GetJSON(url string, cookies ...*http.Cookie) map[string]interface{} {
	return GetJSONUserAgent(url, nil, cookies...)
}

// GetJSONFix downloads contents of url and returns json representation by unmarshal result
func GetJSONFix(url string, cookies ...*http.Cookie) map[string]interface{} {
	userAgent := "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:77.0) Gecko/20100101 Firefox/77.0"
	return GetJSONUserAgent(url, &userAgent, cookies...)
}

// GetJSONUserAgent downloads contents of url and returns json representation by unmarshal result
func GetJSONUserAgent(url string, userAgent *string, cookies ...*http.Cookie) map[string]interface{} {
	req, err := http.NewRequest("GET", url, nil)

	//Fix for SvtPlay
	if userAgent != nil {
		req.Header.Set("User-Agent", *userAgent)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

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
		fmt.Println(url)
		log.Println(err)
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
func GetStringValue(data map[string]interface{}, param string, def *string) *string {
	value := def
	if _, ok := data[param]; ok {
		switch data[param].(type) {
		case string:
			val := data[param].(string)
			value = &val
		case int:
			val := string(data[param].(int))
			value = &val
		case float64:
			if data[param].(float64) == math.Trunc(data[param].(float64)) {
				val := fmt.Sprintf("%.0f", data[param].(float64))
				value = &val
			} else {
				val := fmt.Sprintf("%f", data[param].(float64))
				value = &val
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
				log.Println(err)
				return def
			}
			value = val
		case int:
			value = data[param].(int)
		case float64:
			value = int(data[param].(float64))
		}

	}
	return value
}

//GetFloat64Value from interface
func GetFloat64Value(data map[string]interface{}, param string, def *float64) *float64 {
	value := def
	if _, ok := data[param]; ok {
		switch data[param].(type) {
		case string:
			val, err := strconv.ParseFloat(data[param].(string), 64)
			if err != nil {
				log.Println(err)
				return def
			}
			value = &val
		case int:
			val := float64(data[param].(int))
			value = &val
		case float64:
			val := (data[param].(float64))
			value = &val
		}

	}
	return value
}

//GetBoolValue from interface
func GetBoolValue(data map[string]interface{}, param string, def *bool) *bool {
	value := def
	if _, ok := data[param]; ok {
		switch data[param].(type) {
		case string:
			val, err := strconv.ParseBool(data[param].(string))
			if err != nil {
				log.Println(err)
				return def
			}
			value = &val
		case int:
			val := data[param].(int) != 0
			value = &val
		case float64:
			val := data[param].(float64) != 0
			value = &val
		}

	}
	return value
}

//GetMapValue returns Map from Map
func GetMapValue(data map[string]interface{}, param string) *map[string]interface{} {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if _, ok := data[param]; ok {
		value := data[param].(map[string]interface{})
		return &value
	}
	empty := make(map[string]interface{})
	return &empty
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
