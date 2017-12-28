#!/bin/sh

nohup /app/http-watchmen cron --config /app/conf/app.yml > log.log &
nginx -g "daemon off;"

