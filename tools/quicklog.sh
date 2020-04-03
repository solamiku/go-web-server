#!/bin/bash
#参数说明 1日志目录 2源日志文件 3过滤匹配串 4过滤匹配串接收参数 5日志行数 6临时文件名
#（其中3中可以预留解析参数4的内容，5，6由webserver传入）其余均写在配置中

#echo params 1:$1 2:$2 3:$3 4:$4 5:$5 6:$6 >> params

logdir=$1
mkdir -m 777 $logdir 2> /dev/null

line=500
if [ "$5" -gt 0 ]; then
        line=$5
fi

filter=$(printf "$3" "$4")

#echo $logdir
#echo $line
#echo $filter
#echo $2
#echo $5

grep -E "$filter" "$2" | tail -n $line > $logdir$6
# -P 支持\d -E仅支持[0-9] 但是-P会触发PCRE limits exceeded
#grep -P "$filter" "$2" > $logdir$6
exit 0
