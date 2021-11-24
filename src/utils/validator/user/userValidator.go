package user

import "todoList/src/utils/validator"

type UserValidator struct {
	validator.Validator
}

var UserCreateRules = validator.Rule{"Name": "required;string", "Email": "required", "Phone": "required"}

var SignUpRules = validator.Rule{"Email": "required;email", "Password": "required", "Auth": "required"}

var SignInRules = validator.Rule{"Account": "required", "Auth": "required"}
