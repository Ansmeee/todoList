package mailSVC

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type MailSVC struct{}

func (MailSVC) SendText(subject, content string, receiver ...string) error {
	if len(receiver) == 0 {
		return errors.New("mail send fail: invalid receiver")
	}

	receivers := strings.Join(receiver, ",")

	cmd := exec.Command("bash", "-c", fmt.Sprintf("echo %s | mail -s '%s' '%s' -aFrom:dev@ansme.cc", content, subject, receivers))
	err := cmd.Run()

	if err == nil {
		fmt.Println("send email to: ", receivers)
	}

	return err
}
