#!/bin/bash

# Create a temporary file
tmpfile=$(mktemp /tmp/json.XXXXXX)

# Define the JSON content
json_content='{
  "Title": "DERO-XMR Swap",
  "Content": "Welcome to the DERO-XMR Swap.\n\n
  This swap was made to facilitate private trades.\n\n
  To start, please use the following integrated address to register with the swap:\n\n
  deroi1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqdyvfp4xurnv43hyet5ypkx7an9wvs8jmm4vfz92xvsq93yu4gqvft92xgnxu2fs35t\n\n
  Once you have registered, your DERO address will be paired with a customized XMR address. This \"unique to you\" XMR address will be delivered to your DERO wallet tx history. Please retrieve it there.\n\n
  Afterwards, any time you want to trade XMR for DERO, simply send XMR to the address provided.\n\n
  Currently, the swap only supports XMR->DERO trades."
}'

# Write JSON content to the temporary file
echo "$json_content" > "$tmpfile"

# Make the API call using the temporary file
api_call() {
  curl -X POST -u "user:pass" \
    -H "Content-Type: application/json" \
    -d "@$tmpfile" http://127.0.0.1:3000/api/items
}

# Make the API call
api_call

# Remove the temporary file
rm "$tmpfile"
