package smsSVC

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"todoList/config"
)

type smsSVC struct {
	smsCFG *smsConfig
}

func NewSMSSVC() (*smsSVC, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, err
	}
	smsSVC := new(smsSVC)
	smsSVC.smsCFG = cfg
	return smsSVC, nil
}

func (s *smsSVC) SendCode(code string, mobile ...string) error {
	params := s.generateParams(code, mobile...)

	data, _ := json.Marshal(params)
	res, err := http.Post(s.smsCFG.Host, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("SendCode Error", err)
		return err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	type resp struct {
		Code    int    `form:"code"`
		Message string `form:"message"`
	}

	response := new(resp)
	err = json.Unmarshal(body, response)
	if err != nil {
		log.Println("response parse error", err)
		return err
	}

	if response.Code != 0 {
		log.Println("SendCode Error", string(body))
		return errors.New("SendCode Error")
	}

	return nil
}

func (s *smsSVC) generateParams(code string, mobile ...string) map[string]interface{} {
	params := make(map[string]interface{})
	params["app_id"] = s.smsCFG.AppId
	params["method"] = "sms.message.send"
	params["version"] = "1.0"
	params["timestamp"] = time.Now().UnixMicro()
	params["sign_type"] = "md5"
	params["biz_content"] = s.generateBizCon(code, mobile...)
	params["sign"] = s.generateSign(params)
	return params
}

func (s *smsSVC) generateSign(params map[string]interface{}) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	var signSlice []string
	sort.Strings(keys)
	for _, key := range keys {
		val := params[key]
		signSlice = append(signSlice, fmt.Sprintf("%s=%v", key, val))
	}

	signSlice = append(signSlice, fmt.Sprintf("key=%s", s.smsCFG.AppSecret))
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(signSlice, "&")))))
}

func (s *smsSVC) generateBizCon(code string, mobile ...string) string {
	con := make(map[string]interface{})
	con["mobile"] = mobile
	con["template_id"] = s.smsCFG.TPLId
	con["type"] = 0
	con["params"] = map[string]string{"code": code}
	con["sign"] = s.smsCFG.Sign

	str, _ := json.Marshal(con)
	return string(str)
}

type smsConfig struct {
	Host      string
	AppId     string
	AppSecret string
	TPLId     string
	Sign      string
}

func initConfig() (*smsConfig, error) {
	cfg, err := config.Config()
	if err != nil {
		log.Println("SmsSVC SmsConfig Error", err)
		return nil, err
	}

	host := cfg.Section("sms").Key("host").String()
	appId := cfg.Section("sms").Key("app_id").String()
	appSecret := cfg.Section("sms").Key("app_secret").String()
	tplId := cfg.Section("sms").Key("tpl_id").String()
	sign := cfg.Section("sms").Key("sign").String()
	if host == "" || appId == "" || appSecret == "" || tplId == "" || sign == ""{
		log.Println("SmsSVC SmsConfig Error, invalid config settings")
		return nil, errors.New("invalid config settings")
	}

	conf := new(smsConfig)
	conf.Host = host
	conf.AppId = appId
	conf.AppSecret = appSecret
	conf.Sign = sign
	conf.TPLId = tplId

	return conf, nil
}
