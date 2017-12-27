package mail

import (
	"encoding/json"
	"testing"
)

func TestMailInfo(t *testing.T) {
	data := []byte(`
	{
    "smtp_user":"user",
    "smtp_passwd":"pass",
    "smtp_host":"smtp.exmail.qq.com:25",
    "mailto":[
      "hyhlinux@163.com",
      "2285020853@qq.com"
    ]
  }`)
	var tmp MailInfo
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		t.Errorf("err:%v", err)
	}
	t.Logf("tmp:%v", tmp)

}
func TestSendToMail(t *testing.T) {
	user := "user"
	password := "pass"
	host := "smtp.zoho.com:465"
	to := []string{"to@163.com"}
	subject := "email send by golang"
	out := []string{`2017/04/11 14:44:19 [I] [jobs.go:74] ## STDOUT:`}
	cmd := []string{"pwd", "ls"}
	env := []string{"env1", "env2"}
	err := SendToMail(user, password, host, true, to, subject, out, []string{"err"}, env, cmd, "html")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("send ok")
	}

}
