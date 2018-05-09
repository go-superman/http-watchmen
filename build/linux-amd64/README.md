#### http-watchmen 监控API接口异常状态, 提供命令接口／邮件报警
```yml
- command: [ls, pwd]  #改接口超时后，执行ls;pwd命令.
  cron: '@every 20s'  # 20s 访问一次url
  env: [MONGOPORT=27017]
  mail: null
  name: job-api
  retry_cnt: 3
  retry_time: 5
  # 接口请求超时时间，单位s
  request_timout: 5
  # 接口请求返回异常状态. 默认包含500/404
  request_status: []  # 接口的异常状态
  url: http://47.91.255.0/api/index   #监控的api接口
```

#### http-watchmen 内部功能块
1. http-watchmen 采用cron 语法定义定时任务
2. http-watchmen run 只会执行一次命令，定时任务使用cron参数
3. http-watchmen 接收yml/json配置文件.
4. 可以单独启动该程序，也可以下载到app docker中一起运行.

#### 关于mail配置. smtp协议
```yml
mail:
  mail_subject: '邮件发送'
  mail_to: [tox@163.com]
  smtp_host: smtp.exmail.qq.com:465 #腾讯smtp服务器
  smtp_passwd: xxxxxxx # 邮箱密码
  smtp_tls:
  smtp_user: test@tplinux.net # 邮箱账号
```

#### todo
---
[![Build Status](https://travis-ci.org/go-superman/http-watchmen.svg?branch=master)](https://travis-ci.org/go-superman/http-watchmen)

#### 官方镜像
```bash
docker pull hyhlinux/http-watchmen
```

#### build
```bash
make docker
cp README.md build/linux-amd64
➜  http-watchmen git:(master) ✗ docker build -t hyhlinux/http-watchmen -f Dockerfile.app .
Sending build context to Docker daemon  94.14MB
Step 1/6 : FROM alpine:latest
 ---> 665ffb03bfae
Step 2/6 : MAINTAINER huoyinghui "huoyinghui@apkpure.com"
 ---> Using cache
 ---> 3e9e8ae8871d
Step 3/6 : WORKDIR /app
 ---> Using cache
 ---> 1cc0d7154b39
Step 4/6 : RUN apk update && apk add curl bash tree tzdata     && cp -r -f /usr/share/zoneinfo/Hongkong /etc/localtime
 ---> Using cache
 ---> 364f004e99a3
Step 5/6 : COPY ./build/linux-amd64/ /app/
 ---> 65ab0ddfd132
Removing intermediate container d0ebf2f555b3
Step 6/6 : CMD /app/http-watchmen
 ---> Running in c45d1446da1c
 ---> 7e96c1c1db1f
Removing intermediate container c45d1446da1c
Successfully built 7e96c1c1db1f
Successfully tagged hyhlinux/http-watchmen:latest
```