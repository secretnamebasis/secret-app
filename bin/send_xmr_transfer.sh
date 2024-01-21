curl http://127.0.0.1:18082/json_rpc \
  -d '{
    "jsonrpc": "2.0",
    "id": "0",
    "method": "transfer",
    "params": {
      "destinations": [
        {
          "amount": 100000000000,
          "address": "7BnERTpvL5MbCLtj5n9No7J5oE5hHiB3tVCK5cjSvCsYWD2WRJLFuWeKTLiXo5QJqt2ZwUaLy2Vh1Ad51K7FNgqcHgjW85o"
        },
        {
          "amount": 200000000000,
          "address": "75sNpRwUtekcJGejMuLSGA71QFuK1qcCVLZnYRTfQLgFU5nJ7xiAHtR5ihioS53KMe8pBhH61moraZHyLoG4G7fMER8xkNv"
        }
      ],
      "account_index": 0,
      "subaddr_indices": [0],
      "priority": 0,
      "ring_size": 7,
      "get_tx_key": true
    }
  }' -H 'Content-Type: application/json'
