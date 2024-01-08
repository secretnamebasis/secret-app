user="secretnamebasis"
monero_pass="bargraph-chivalry-bullhorn"
ip="192.168.12.176"
port="28088"
n=2953189
min=$(($n-1))
curl -s -X POST http://$ip:$port/json_rpc \
-u $user:$monero_pass --digest \
-H 'Content-Type: application/json' \
-d '{
  "jsonrpc":"2.0",
  "id":"0",
  "method":"get_transfers",
  "params":
  {"in":true,
  "filter_by_height":true,
  "min_height":'"$min"',
  "max_height":'"$n"'
  }
  }'
