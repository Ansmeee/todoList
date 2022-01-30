package captcha

import (
	"github.com/dchest/captcha"
	"net/http"
	"strings"
)

type CaptchaService struct {}

func (CaptchaService) GenerateID() string {
	captchaID := captcha.NewLen(6)
	return captchaID
}

func (CaptchaService) GenerateImg(writer http.ResponseWriter, source string) {
	sourceMap := strings.Split(source, ".")
	id := sourceMap[0]
	writer.Header().Set("Content-Type", "image/png")
	captcha.WriteImage(writer, id, captcha.StdWidth, captcha.StdHeight)
	return
}

func (CaptchaService) Verify(id, value string) bool{
	return captcha.VerifyString(id, value)
}