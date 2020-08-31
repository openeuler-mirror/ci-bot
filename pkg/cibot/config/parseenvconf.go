package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// ParseEnvConf  Parse environment variables according to tag.
// only supports string、bool、int、float type variable parsing settings.
func ParseEnvConf(obj interface{}, tagName string) (err error){
	defer func() {
		v := recover()
		if v != nil {
			err = errors.New(fmt.Sprintf("%v",v))
		}
	}()
	if obj == nil {
		return errors.New("obj param can not be nil")
	}
	if tagName == "" {
		tagName = "envVariable"
	}
	v := reflect.ValueOf(obj)
	if e := v.Type().Elem(); e.Kind() != reflect.Struct {
		return errors.New("obj param only support  struct")
	}
	st := v.Type().Elem()
	stFiledNum := st.NumField()
	for i := 0; i < stFiledNum; i++ {
		f := st.Field(i)
		ft := f.Tag.Get(tagName)
		if ft == "" || ft == "-" {
			continue
		}
		environ := os.Getenv(ft)
		newValue := v.Elem().FieldByName(f.Name)
		switch f.Type.Kind() {
		case reflect.String:
			if newValue.CanSet() && environ != "" {
				newValue.SetString(environ)
			}
		case reflect.Int:
			if newValue.CanSet() && environ != "" {
				val, e := strconv.Atoi(environ)
				if e == nil {
					newValue.SetInt(int64(val))
				}else {
					err = e
				}
			}
		case reflect.Bool:
			if newValue.CanSet() && environ != "" {
				val, e := strconv.ParseBool(environ)
				if e == nil {
					newValue.SetBool(val)
				}else {
					err = e
				}
			}
		case reflect.Float32:
			if newValue.CanSet() && environ != "" {
				val, e := strconv.ParseFloat(environ,32)
				if e == nil {
					newValue.SetFloat(val)
				}else {
					err = e
				}
			}
		case reflect.Float64:
			if newValue.CanSet() && environ != "" {
				val, e := strconv.ParseFloat(environ,64)
				if e == nil {
					newValue.SetFloat(val)
				}else {
					err = e
				}
			}
		}
	}
	return
}