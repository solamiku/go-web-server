#!/bin/bash
#参数说明 1日志目录 2源日志文件 3过滤匹配串 4过滤匹配串接收参数 5日志行数 6临时文件名
#（其中3中可以预留解析参数4的内容，5，6由webserver传入）其余均写在配置中

logdir=$1
mkdir $logdir 2> /dev/null

line=20000
if [ "$5" -gt 0 ]; then
        line=$5
fi

filter=$(printf "$3" "$4")

#echo $logdir
#echo $line
#echo $filter
#echo $2
#echo $5

grep $filter $2 | tail -n $line > $logdir$6
exit 0
