#!/bin/bash
query="__tag__:__path__:\"/home/stararc/$1/log/gamesrv.err.log\""
ftime="$2"
ttime="$3"
output="$4"
echo try get query:$query from:$ftime to:$time output:$output
prj="stararc"
lstore="stararc-logstore"
endp="cn-shanghai.log.aliyuncs.com"
jmes="[*].join(' ',[to_string(@.time),to_string(@.level),to_string(@.message)])| join('\n', map(&to_string(@), @))"
aid=""
akey=""

aliyunlog log get_log_all --project="$prj" --logstore="$lstore" --query="$query" --from_time="$ftime" --to_time="$ttime" --region-endpoint="$endp" --format-output=no_escape --jmes-filter="$jmes" --access-id="$aid" --access-key="$akey" > $output
