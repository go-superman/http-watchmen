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
	"fmt"
	"github.com/go-superman/http-watchmen/job"
	"github.com/go-superman/http-watchmen/logger"

	"encoding/json"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var jobName string
var redisAddr string
var redisPasswd string
var redisDBIndex int
var serPort int
var logLevel int

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run job",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := RootCmd.Flag("config")
		err := checkFlag(config)
		if err != nil {
			logger.Errorf("err:%v", err)
			return
		}

		bcJob, err := job.LoadConf(config.Value.String())
		if err != nil {
			logger.Errorf("err:%v config:%v", err, bcJob)
			return
		}
		b, err := json.Marshal(bcJob)
		if err != nil {
			logger.Errorf("err:%v config:%v", err, bcJob)
			return
		}

		logger.Infof("config-json:%v", string(b))
		jobName := cmd.Flag("jobname")
		err = checkFlag(jobName)
		if err != nil {
			logger.Errorf("err:%v", err)
			return
		}

		for _, job := range bcJob.Jobs {
			logger.Debugf("job:%v job.name:%v   jobname:%v  bcJob.Mail:%v",
				job, job.Name, jobName.Value.String(), bcJob.Mail)
			if jobName.Value.String() == job.Name {
				job.Mail = &bcJob.Mail
				job.Run()
			}
		}

	},
}

func checkFlag(f *flag.Flag) (err error) {
	if f == nil {
		err = fmt.Errorf("%v..%v", f.Name, f.Usage)
		return
	}
	//
	if f != nil && f.Value.String() == "" {
		err = fmt.Errorf("%v can not be empty..%v", f.Name, f.Usage)
		return
	}

	return
}

func init() {
	RootCmd.PersistentFlags().StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis ser addr")
	RootCmd.PersistentFlags().StringVar(&redisPasswd, "redisPasswd", "", "redis ser passwd")
	RootCmd.PersistentFlags().IntVar(&redisDBIndex, "redisDBIndex", 0, "redis ser DB 0/1/..")
	RootCmd.PersistentFlags().IntVar(&logLevel, "logLevel",  6, "info")

	RootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVar(&jobName, "jobName", "", "which job you want to run")
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
