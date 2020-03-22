package validation

import (
    "reflect"

    "github.com/go-playground/validator"
)

//IsSlice check if field kind is equal to slice
func IsSlice(fl validator.FieldLevel) bool {
    if fl.Top().Kind() == reflect.Slice {
        return true
    }
    return false
}

//IsStringElem check if field element kind is equal to string
func IsStringElem(fl validator.FieldLevel) bool {
    t := fl.Top().Type()
    if t.Elem().Kind() == reflect.String {
        return true
    }
    return false
}
