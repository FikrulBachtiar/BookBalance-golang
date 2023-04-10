package model

import "github.com/go-playground/validator/v10"

type PayloadValidator struct {
	Validator *validator.Validate
}

func (pv *PayloadValidator) Validate(i interface{}) error {
	return pv.Validator.Struct(i)
}