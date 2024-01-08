#! usr/env/bin
curl -X POST \
  http://192.168.12.208:10104/json_rpc \
  -u secret:pass \
  -H 'content-type: application/json' \
  -d '{
    "jsonrpc": "2.0",
    "id": "1",
    "method": "GetTransfers",
    "params": {
        "out": true,
        "in": true
    }
}'
