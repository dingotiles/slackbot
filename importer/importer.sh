#!/bin/bash
touch $1
> $1
robots=(
    "github.com/dingotiles/slackbot/robots/downloads3"
    "github.com/dingotiles/slackbot/robots/downloads3rc"
    "github.com/dingotiles/slackbot/robots/downloadkafka"
    "github.com/dingotiles/slackbot/robots/downloadkafka"
    "github.com/dingotiles/slackbot/robots/downloadredis"
    "github.com/dingotiles/slackbot/robots/downloadpostgresql"
    "github.com/dingotiles/slackbot/robots/downloadpostgresqlrc"
    "github.com/dingotiles/slackbot/robots/downloadprometheus"
    "github.com/dingotiles/slackbot/robots/downloadprometheusrc"
    "github.com/dingotiles/slackbot/robots/downloadsecrets"
    "github.com/dingotiles/slackbot/robots/downloadshield"
    "github.com/dingotiles/slackbot/robots/downloadshieldrc"
)

echo "package importer

import (" >> $1

for robot in "${robots[@]}"
do
    echo "    _ \"$robot\" // automatically generated import to register bot, do not change" >> $1
done
echo ")" >> $1

gofmt -w -s $1
