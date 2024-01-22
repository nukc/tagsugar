package tagsugar

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	Http    = ""
	Debug   = false
	hostMap = make(map[string]string, 0)
)

func AddHost(key string, host string) {
	hostMap[key] = host
}

func Lick(data interface{}) {
	v := reflect.ValueOf(data)
	k := v.Kind()
	resolveValue(v, k)
}

// get the value that the Elem and the Elem's Kind.
func getEkByValue(value reflect.Value) (v reflect.Value, k reflect.Kind) {
	v = value.Elem()
	k = v.Kind()
	return
}

func resolveValue(v reflect.Value, k reflect.Kind) {
	switch k {
	case reflect.Ptr:
		v, k = getEkByValue(v)
		resolveValue(v, k)
	case reflect.Slice:
		arraySlice(v)
	case reflect.Struct:
		resolveField(v)
	case reflect.Interface:
		v, k = getEkByValue(v)
		resolveValue(v, k)
	default:
		if Debug {
			log.Print("Ignore kind: " + k.String())
		}
	}
}

// slice
func arraySlice(v reflect.Value) {
	count := v.Len()
	for i := 0; i < count; i++ {
		item := v.Index(i)
		k := item.Kind()
		switch k {
		case reflect.Interface:
			resolveValue(item, k)
		case reflect.Struct, reflect.Array, reflect.Slice:
			resolveField(item)
		case reflect.Ptr:
			value, _ := getEkByValue(item)
			resolveField(value)
		}
	}
}

// resolve the value that field
func resolveField(value reflect.Value) {
	if !value.IsValid() {
		return
	}
	p := value.Type()
	l := p.NumField()
	for i := 0; i < l; i++ {
		field := value.Field(i)
		k := field.Kind()
		if k == reflect.Slice {
			arraySlice(field)
			continue
		} else if k == reflect.Ptr {
			field, k = getEkByValue(field)
			resolveValue(field, k)
			continue
		} else if k == reflect.Interface {
			resolveValue(field, k)
			continue
		}

		sf := p.Field(i)
		options := parseTag(sf.Tag.Get("ts"))
		if err := changeField(value, field, options); err != nil && Debug {
			log.Print(err)
		}

	}
}

// change field according to tag options
func changeField(v reflect.Value, field reflect.Value, options tagOptions) error {
	if _, ok := options["-"]; ok || len(options) == 0 {
		return nil
	}

	if _, ok := options["initial"]; ok {
		// new initial value
		obj := newInterface(field)
		field.Set(reflect.ValueOf(obj).Elem())
		return nil
	}

	url := options["url"]
	if url == "http" {
		if field.CanSet() {
			var s = field.String()
			if s != "" && !strings.HasPrefix(s, "http") {
				field.Set(reflect.ValueOf(Http + s))
			}
		}
	}

	host := options["host"]
	if host != "" {
		if field.CanSet() {
			s := field.String()
			value := hostMap[host]
			if s != "" && value != "" {
				field.Set(reflect.ValueOf(value + s))
			}
		}
	}

	cName := options["assign_to"]
	if cName != "" {
		cFiled := v.FieldByName(cName)
		if !cFiled.IsValid() {
			return errors.New("The field that needs to be assigned does not exist, ts tag: copyTo(" + cName + ")")
		} else if !cFiled.CanSet() {
			return errors.New("The copy to field can not set ")
		}

		switch options["assign_type"] {
		case "unmarshal":
			str := field.String()
			if str == "" {
				return errors.New("unexpected end of JSON input")
			}

			obj := newInterface(cFiled)
			err := json.Unmarshal([]byte(str), &obj)
			if err == nil {
				ov := reflect.ValueOf(obj)
				cFiled.Set(ov.Elem())
			} else {
				return err
			}
		case "bool":
			b := assignBool(field)
			cFiled.SetBool(b)
		case "raw":
			cFiled.Set(field)
		case "string":
			s := assignString(field)
			cFiled.SetString(s)
		default:
			cFiled.Set(field)
		}
	}

	return nil
}

// returns interface{}
func newInterface(field reflect.Value) interface{} {
	typ := reflect.Indirect(field).Type()
	obj := reflect.New(typ).Interface()
	return obj
}

func assignBool(v reflect.Value) bool {
	k := v.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 1
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 1
	case reflect.String:
		return v.String() == "1"
	case reflect.Bool:
		return v.Bool()
	}
	panic(&reflect.ValueError{Method: "ts tag assignBool", Kind: k})
}

func assignString(v reflect.Value) string {
	k := v.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	}
	panic(&reflect.ValueError{Method: "ts tag assignString", Kind: k})
}
