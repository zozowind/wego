package util

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

// URLValueToXML trans url.Values to xml string
func URLValueToXML(params url.Values) string {
	resultXML := "<xml>"
	for key, values := range params {
		value := values[0]
		resultXML += "<" + key + "><![CDATA[" + value + "]]></" + key + ">"
	}
	resultXML += "</xml>"
	return resultXML
}

// StructToURLValue trans struct to url.Values
func StructToURLValue(obj interface{}, t string) (url.Values, error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	if !(objT.Kind() == reflect.Ptr && objT.Elem().Kind() == reflect.Struct) {
		return nil, fmt.Errorf("%v must be a struct pointer", obj)
	}
	objT = objT.Elem()
	objV = objV.Elem()
	params := url.Values{}
	for i := 0; i < objT.NumField(); i++ {
		fieldV := objV.Field(i)
		if !fieldV.CanSet() {
			continue
		}
		fieldT := objT.Field(i)
		tag := fieldT.Tag.Get(t)
		if tag == "-" {
			continue
		} else if tag == "" {
			tag = fieldT.Name
		}
		value := ""
		switch fieldT.Type.Kind() {
		case reflect.Bool:
			value = "false"
			if fieldV.Bool() {
				value = "true"
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = fmt.Sprintf("%d", fieldV.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = fmt.Sprintf("%d", fieldV.Uint())
		case reflect.Float32, reflect.Float64:
			value = fmt.Sprintf("%f", fieldV.Float())
			re := regexp.MustCompile(`[\.]?0+$`)
			value = re.ReplaceAllString(value, "")
		case reflect.String:
			value = fieldV.String()

		}
		if value != "" {
			params.Add(tag, value)
		}
	}
	return params, nil
}
