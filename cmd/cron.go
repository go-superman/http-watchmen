// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"github.com/go-superman/http-watchmen/job"
	"github.com/go-superman/http-watchmen/logger"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"
)

const ONEHOUR = 60 * time.Minute

// cronCmd represents the cron command
var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Run cron backup job",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: cronCall,
}

func init() {
	RootCmd.AddCommand(cronCmd)
	//RootCmd.PersistentFlags().String("config", "", "Enter a conf filename")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cronCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cronCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func cronCall(cmd *cobra.Command, args []string) {
	// TODO run cron backup jobs
	config := RootCmd.Flag("config")
	err := checkFlag(config)
	if err != nil {
		logger.Errorf("err:%v", err)
		return
	}

	bcJob, err := job.LoadConf(config.Value.String())
	if err != nil {
		logger.Errorf("err:%v tmpyml:%v", err, bcJob)
		return
	}

	b, err := json.Marshal(bcJob)
	if err != nil {
		logger.Errorf("err:%v config:%v", err, bcJob)
		return
	}

	logger.Infof("config-json:%v", string(b))
	location, err := time.LoadLocation(bcJob.Timezone)
	if err != nil {
		logger.Errorf("ERROR : %s", err)
		return
	}
	resetLog()
	c := cron.NewWithLocation(location)
	for index, task := range bcJob.Jobs {
		logger.Debugf("index:%v job:%v bcjob.mail:%v", index, task, bcJob.Mail)
		task.Mail = &bcJob.Mail
		task.Timezone = bcJob.Timezone
		task.RedisAddr = parseRootCmd("redisAddr")
		task.RedisPasswd = parseRootCmd("redisPasswd")
		task.RedisKeyPrefix = parseRootCmd("redisKeyPrefix")
		task.RedisDB = parseRootCmdInt("redisDBIndex", 0)
		c.AddJob(task.Cron, task)
	}

	//启动启动定时任务
	c.Start()
	displayJob(c)
	defer c.Stop()
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	func() {
		for {
			select {
			//case <-time.After(ONEHOUR):
			//	displayJob(c)
			case <-sig:
				logger.Debugf("cron done")
				os.Exit(0)
				return
			}
		}
	}()
	return
}

func displayJob(c *cron.Cron) {
	if c == nil {
		logger.Debugf("no job")
		return
	}

	for index, job := range c.Entries() {
		logger.Debugf("cur_id:%v job-next-time:%v", index, job.Next)
	}
	return
}

func parseRootCmd(name string) (value string){
	flag := RootCmd.Flag(name)
	err := checkFlag(flag)
	if err != nil {
		logger.Debugf("err:%v, will be empty", err)
		value = flag.DefValue
	}else{
		value = flag.Value.String()
	}
	return value
}

//value: 当解析出错是，期望的值
func parseRootCmdInt(name string, value int)  (ret int){
	flagString := parseRootCmd(name)
	ret, err := strconv.Atoi(flagString)
	if err != nil {
		logger.Warnf("err:%v", err)
		ret = value
	}
	return ret
}

func resetLog()  {
	level := parseRootCmdInt("logLevel", 7)
	logger.ResetLevel(level)
}