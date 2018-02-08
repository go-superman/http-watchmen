package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-superman/http-watchmen/logger"
	"html/template"
	"net"
	"net/smtp"
	"strings"
)

var mailTpl *template.Template

func init() {
	mailTpl, _ = template.New("mail_tpl").Parse(`
	你好 {{.username}}，<br/>

<p>任务详情：</p>
<p>
命令数量:{{.cmdlen}}
<br>
{{range $index, $elem := .cmds}}
    {{$index}}:{{$elem}}
    <br>
{{end}}
<br>
环境变量:
{{range $index, $elem := .env}}
    {{$index}}:{{$elem}}
    <br>
{{end}}
</p>
<p>-------------以下是任务执行输出-------------</p>
<p>
{{range $index, $elem := .output}}
    {{$index}}:{{$elem}}
    <br>
{{end}}
</p>
<p>-------------以下是任务执行错误信息-------------</p>
<p>
{{range $index, $elem := .outputerr}}
    {{$index}}:{{$elem}}
    <br>
{{end}}
</p>
<p>
--------------------------------------------<br />
本邮件由系统自动发出，请勿回复<br />
</p>
`)

}

type MailInfo struct {
	MailUser      string   `yaml:"smtp_user" json:"smtp_user"`
	MailPasswd    string   `yaml:"smtp_passwd" json:"smtp_passwd"`
	MailHost      string   `yaml:"smtp_host" json:"smtp_host"`
	MailEnableTls bool     `yaml:"smtp_tls" json:"smtp_tls"`
	MailTo        []string `yaml:"mail_to" json:"mail_to"`
	MailSubject   string   `yaml:"mail_subject" json:"mail_subject"`
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	// TLS config
	host, _, _ := net.SplitHostPort(addr)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		logger.Errorf("Dialing Error:%v", err)
		return nil, err
	}
	//分解主机端口字符串
	return smtp.NewClient(conn, host)
}

func SendMailUsingTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	c, err := Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				logger.Errorf("Error during AUTH:%v", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func SendMail(mail *MailInfo, out []string, outErr []string, env, cmd []string, mailtype string) error {
	if mail == nil {
		return errors.New("mail is nil!!")
	}
	return SendToMail(mail.MailUser, mail.MailPasswd, mail.MailHost, mail.MailEnableTls, mail.MailTo, mail.MailSubject,
		out, outErr, env, cmd, mailtype)
}
func SendToMail(user, password, host string, enableTls bool, to []string, subject string, out, outErr []string, env, cmd []string, mailtype string) (err error) {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	data := make(map[string]interface{})
	data["username"] = to
	data["output"] = out
	data["outputerr"] = outErr
	data["cmds"] = cmd
	data["env"] = env
	data["cmdlen"] = len(cmd)
	content := new(bytes.Buffer)
	mailTpl.Execute(content, data)
	body := content.String()
	//logger.Debugf("body:%v", body)
	msg := []byte("To: " + strings.Join(to, ";") + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	if !enableTls {
		err = smtp.SendMail(host, auth, user, to, msg)
	} else {
		err = SendMailUsingTLS(host, auth, user, to, msg)
	}
	if err != nil {
		err = fmt.Errorf("err:%v to:%v host:%v auth:%v user:%v", err, to, host, auth, user)
	}
	return err
}
