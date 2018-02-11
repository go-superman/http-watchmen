package job

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-superman/http-watchmen/logger"
	"github.com/go-superman/http-watchmen/mail"
	"github.com/go-superman/http-watchmen/utils"
	"github.com/go-superman/http-watchmen/storage"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
	"io"
	iot "io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

const APPBACKUPPATH = "/tmp/smart_backup"

type JobConfig struct {
	Timezone string        `yaml:"timezone" json:"timezone"` //可以自定义cron的执行时区
	Jobs     []*Job        `yaml:"jobs,flow" json:"jobs"`
	Mail     mail.MailInfo `yaml:"mail" json:"mail"` // 邮件
}

type Job struct {
	Name      string         `yaml:"name" json:"name"`
	Url       string         `yaml:"url" json:"url"`             // url 健康检查
	RetryCnt  int            `yaml:"retry_cnt" json:"retry_cnt"` // 重试次数
	RetryTime int            `yaml:"" json:"retry_time"`
	Timezone  string         `yaml:"-" json:"-"`             // 无需设置，程序会把BackupJobConfig.Timezone复制过来
	Cron      string         `yaml:"cron" json:"cron"`       // crontab
	Command   []string       `yaml:"command" json:"command"` // shell command
	ENV       []string       `yaml:"env" json:"env"`         // shell env
	Mail      *mail.MailInfo `yaml:"mail" json:"mail"`       // 邮件
	RedisAddr string	   `yaml:"redis_addr" json:"redis_addr"`
	RedisPasswd string	   `yaml:"redis_passwd" json:"redis_passwd"`
	RedisKeyPrefix string	   `yaml:"redis_key_prefix" json:"redis_key_prefix""`
	RedisDB  int		   `yaml:"redis_db" json:"redis_db"`
}

func (job *Job) Run() {
	var err error
	var data string
	var allTextOut []string
	var allTextErr []string
	outPutPath := path.Join(APPBACKUPPATH, job.Name)

	logger.Debugf("job.name:%v start .. ", job.Name)
	defer displayNext(job.Name, job.Cron, job.Timezone)
	data, err = utils.HealthCheck(job.Url, job.RetryCnt, time.Duration(job.RetryTime)*time.Second)
	defer job.saveData(data)
	if err == nil {
		// 本次健康检查成功
		return
	}

	// err不为空时，执行用户指定的cmd
	osEnv := os.Environ()
	for index, cmdTmp := range job.Command {
		allTextOut = append(allTextOut, fmt.Sprintf("CMD_INDEX: %d", index))
		allTextErr = append(allTextErr, fmt.Sprintf("CMD_INDEX: %d", index))
		var cmd *exec.Cmd
		logger.Debugf("CMD##:%v", cmdTmp)
		cmd = exec.Command("sh", "-c", cmdTmp)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			logger.Errorf("err:%v", err)
			continue
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			logger.Errorf("err:%v", err)
			continue
		}
		cmd.Env = append(osEnv, job.ENV...)
		cmd.Env = append(cmd.Env, fmt.Sprintf("OUTPUTPATH=%s", outPutPath))
		logger.Debugf("ENV##:%v...", cmd.Env)
		var w sync.WaitGroup
		go func(r io.Reader) {
			w.Add(1)
			defer w.Done()
			allTextOut = append(allTextOut, ReadDataFromStdout(r)...)
		}(stdout)

		go func(r io.Reader) {
			w.Add(1)
			defer w.Done()
			allTextErr = append(allTextErr, ReadDataFromStderr(r)...)
		}(stderr)

		err = cmd.Run()
		if err != nil {
			logger.Errorf("ERR##:%v index:%v  cmd-path:%v cmd-arg:%v ", err, index, cmd.Path, cmd.Args)
			goto SENDMAIL
		}

		w.Wait()
	}

SENDMAIL:
	// 用户命令执行失败时，发送邮件通知
	logger.Errorf("err:%v will send email", err)
	if job.Mail == nil {
		logger.Errorf("JobMail: %v", job.Mail)
		return
	}
	job.Mail.MailEnableTls = true
	allTextErr = append(allTextErr, fmt.Sprintf("%v", err))
	allTextOut = append(allTextOut, fmt.Sprintf("%v", job.Name))
	tmpMail := mail.MailInfo{
		MailUser: job.Mail.MailUser,
		MailPasswd: job.Mail.MailPasswd,
		MailHost: job.Mail.MailHost,
		MailSubject: fmt.Sprintf("%v-%v:%v", job.Mail.MailSubject, job.RedisKeyPrefix, job.Name),
		MailTo: job.Mail.MailTo,
		MailEnableTls: job.Mail.MailEnableTls,
	}
	err = mail.SendMail(&tmpMail, allTextOut, allTextErr, job.ENV, job.Command, "html")
	if err != nil {
		logger.Errorf("send mail err:%v tmpMail.Subject:%v", err, tmpMail.MailSubject)
	}
	return
}

func ReadDataFromStdout(r io.Reader) (data []string) {
	buf := bufio.NewScanner(r)
	for buf.Scan() {
		tmpData := buf.Text() + "\n"
		io.Copy(os.Stdout, strings.NewReader(tmpData))
		data = append(data, tmpData)
	}
	return
}
func ReadDataFromStderr(r io.Reader) (data []string) {
	buf := bufio.NewScanner(r)
	for buf.Scan() {
		tmpData := buf.Text() + "\n"
		io.Copy(os.Stderr, strings.NewReader(tmpData))
		data = append(data, tmpData)
	}
	return
}

func displayNext(name, spec string, tz string) {
	schedule, err := cron.Parse(spec)
	if err != nil {
		logger.Errorf("err:%v scheule:%v spec:%v", err, schedule, spec)
		return
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		logger.Errorf("err: %v  tz:%v", err, tz)
		return
	}
	next := schedule.Next(time.Now().In(location))
	logger.Debugf("%v: nextRun:%v", name, next)
	return
}

func LoadConf(filepath string) (tmpBackupJobConfig *JobConfig, err error) {
	if filepath == "" {
		return nil, errors.New("filepath is empty, must use --config xxx.yml/json")
	}

	data, err := iot.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	tmpBackupJobConfig = &JobConfig{}

	if strings.HasSuffix(filepath, ".json") {
		err = json.Unmarshal(data, tmpBackupJobConfig)
	} else if strings.HasSuffix(filepath, ".yml") || strings.HasSuffix(filepath, ".yaml") {
		err = yaml.Unmarshal(data, tmpBackupJobConfig)
	} else {
		return nil, errors.New("you config file must be json/yml")
	}

	if err != nil {
		return nil, err
	}

	return tmpBackupJobConfig, nil
}
func (job *Job) saveData(data string) {
	key := fmt.Sprintf("%v:%v", job.RedisKeyPrefix, job.Name)
	client := storage.NewClient(job.RedisAddr, job.RedisPasswd, job.RedisDB)
	err := client.Set(key, data, 0).Err()
	if err != nil {
		logger.Errorf("saveData:%v err:%v", data, err)
	}
}

