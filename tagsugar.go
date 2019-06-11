package tagsugar

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
)

var (
	Http = ""
)

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
		break
	case reflect.Slice:
		arraySlice(v)
		break
	case reflect.Struct:
		resolveField(v)
		break
	case reflect.Interface:
		v, k = getEkByValue(v)
		resolveValue(v, k)
	default:
		log.Print("Ignore kind: " + k.String())
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
			break
		case reflect.Struct, reflect.Array, reflect.Ptr, reflect.Slice:
			resolveField(item)
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
			resolveField(field)
			continue
		} else if k == reflect.Interface {
			resolveValue(field, k)
			continue
		}

		sf := p.Field(i)
		options := parseTag(sf.Tag.Get("ts"))
		if err := changeField(value, field, options); err != nil {
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
			if !strings.HasPrefix(s, "http") {
				field.Set(reflect.ValueOf(Http + s))
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
			break
		case "bool":
			b := assignBool(field)
			cFiled.SetBool(b)
		case "raw":
			cFiled.Set(field)
			break
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
	case reflect.String:
		return v.String() == "1"
	case reflect.Bool:
		return v.Bool()
	}
	panic(&reflect.ValueError{"ts tag assignBool", k})
}
