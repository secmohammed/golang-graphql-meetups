package validation

import (
    "reflect"

    "github.com/go-playground/validator"
)

//IsSlice check if field kind is equal to slice
func IsSlice(fl validator.FieldLevel) bool {
    if fl.Field().Kind() == reflect.Slice {
        return true
    }
    return false
}

//IsStringElem check if field element kind is equal to string
func IsStringElem(fl validator.FieldLevel) bool {
    for i := 0; i < fl.Field().Cap(); i++ {
        if reflect.TypeOf(fl.Field().Index(i).Interface().(string)).Kind() != reflect.String {
            return false
        }
    }
    return true
}
