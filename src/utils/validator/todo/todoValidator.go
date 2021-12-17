package todo

import "todoList/src/utils/validator"

type TodoValidator struct {
	validator.Validator
}

var TodoCreateRules = validator.Rule{}