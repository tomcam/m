package app

import (
	"errors"
	"fmt"
	"reflect"
)

// fieldIsStringType() determines whether the struct passed in the
// argument has a field named by key that is of type string.
func fieldIsStringType(obj interface{}, key string) bool {
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		return false
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	for i := 0; i < t.NumField(); i++ {
		fieldType := fmt.Sprint(v.Field(i).Type())
		if fieldType == "string" && v.Type().Field(i).Name == key {
			return true
		}
	}
	return false
}

// setField writes a value in a struct (well, interface)
// but returns an error if there is no field by
// that name in the struct, or if it's read only.
// Based on https://stackoverflow.com/a/26746461
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

//	for k, v := range app.Page.frontMatterRaw {
//xx		setFieldMust(&app.Page.FrontMatter, k, v)

// setField writes a value in a struct (well, interface)
// but returns an error if there is no field by
// that name in the struct, or if it's read only.
// Based on https://stackoverflow.com/a/26746461

// setFieldMust() writes a value from obj with the field named
// name, to the field by that same name in the interface named value.
// It is identical to setField but strips the
// error checking.
// Use in cases like frontMatterRawtoStruct()
// where the structure type is known in advance.
func setFieldMust(obj interface{}, name string, value interface{}) {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return
	}

	if !structFieldValue.CanSet() {
		return
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return
	}

	structFieldValue.Set(val)
	return
}

// structFieldByNameStrMust() takes any struct and field name (as a string)
// passed in at runtime and returns the string value of that field.
// It returns an empty string if the
// object passed in isn't a struct, or if the named field isn't a struct.
func structFieldByNameStrMust(obj interface{}, field string) string {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		return ""
	}
	/*
		kind := v.FieldByName(field).Kind()
		if kind != reflect.String {
			return ""
		}
	*/
	return (fmt.Sprint(v.FieldByName(field)))
}

// structHasField() returns true if a struct passed to it at runtime contains a field name passed as a string
func structHasField(obj interface{}, field string) bool {
	v := reflect.ValueOf(obj)
	if reflect.TypeOf(obj).Kind() != reflect.Struct {
		return false
	}
	return reflect.Indirect(v).FieldByName(field).IsValid()
}
