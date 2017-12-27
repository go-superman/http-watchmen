#### http-watchmen
1. http-watchmen 采用cron 语法定义定时任务
2. http-watchmen run 只会执行一次命令，定时任务使用cron参数
3. http-watchmen 接收yml/json配置文件.
4. 可以单独启动该程序，也可以下载到app docker中一起运行.

#### todo
1.任务状态保存到redis
---
[![Build Status](https://travis-ci.org/go-superman/http-watchmen.svg?branch=master)](https://travis-ci.org/go-superman/http-watchmen)
