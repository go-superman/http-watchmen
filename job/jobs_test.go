package job

import (
	"github.com/go-superman/http-watchmen/mail"
	"gopkg.in/yaml.v2"
	iot "io/ioutil"
	"testing"
)

func TestNewConf(t *testing.T) {
	mailInfo := mail.MailInfo{
		MailUser:      "user",
		MailPasswd:    "pass",
		MailHost:      "smtp.zoho.com:465",
		MailEnableTls: true,
		MailTo: []string{
			"huoyinghui@apkpure.net",
		},
	}
	jobs := []*Job{
		&Job{
			Name:      "job-api",
			Url:       "http://47.91.255.0/api/index",
			RetryCnt:  3,
			RetryTime: 5,
			Cron:      "@every 20s",
			Command: []string{
				"ls",
				"pwd",
			},
			ENV: []string{
				"MONGOPORT=27017",
			},
		},
	}
	conf := JobConfig{
		Timezone: "Asia/Shanghai",
		Jobs:     jobs,
		Mail:     mailInfo,
	}

	data, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%v", string(data))
	if err := iot.WriteFile("../conf/app_template.yml", data, 0644); err != nil {
		t.Fatal(err)
	}

}

func TestLoadConf(t *testing.T) {
	bcJobYml, err := LoadConf("../conf/app_template.yml")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("bcJob2:%v", bcJobYml)
}
