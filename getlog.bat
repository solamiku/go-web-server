@echo off
set QUERY="%1"
set OUTPUT="%4"
echo try get uid:%QUERY% from:%2 to:%3 output:%OUTPUT%

aliyunlog log get_log_all --project="stararc" --logstore="stararc-logstore" --query=%QUERY% --from_time=%2 --to_time=%3 --region-endpoint="cn-shanghai.log.aliyuncs.com" --format-output=no_escape --jmes-filter="[*].join(' ',[to_string(@.time),to_string(@.level),to_string(@.user),to_string(@.message)])| join('\n', map(&to_string(@), @))" --access-id="" --access-key="" 1>%OUTPUT%
