package list

import "todoList/src/utils/validator"

type ListValidator struct {
	validator.Validator
}

var CreateRules = validator.Rule{"Title": "required"}