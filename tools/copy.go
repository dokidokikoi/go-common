package tools

import (
	"errors"
	"reflect"
)

func PropertiesCopy(target, source any) error {
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}
	tv = tv.Elem()
	tt := reflect.TypeOf(tv.Interface())

	sv := reflect.ValueOf(source)
	if sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
		source = sv.Interface()
	}
	st := reflect.TypeOf(source)

	for i := sv.NumField() - 1; i >= 0; i-- {
		tf := tv.FieldByName(st.Field(i).Name)
		if tf.IsZero() {
			if tf.CanSet() && tt.Field(i).Type.Kind() == st.Field(i).Type.Kind() {
				tf.Set(sv.Field(i))
			}
		}
	}
	return nil
}
