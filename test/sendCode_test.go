package test

import (
	"testing"
	"todoList/config"
	"todoList/src/services/smsSVC"
)

func TestSendCode(t *testing.T)  {
	config.InitConfig()

	smsSVC := smsSVC.NewSMSSVC()
	m := []string{"15692237913"}
	smsSVC.SendCode("5689", m...)
}