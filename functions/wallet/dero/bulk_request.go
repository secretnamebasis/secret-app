package dero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
)

// MakeBulkTransfer performs a bulk transfer of DERO coins.
func MakeBulkTransfer(transfers []rpc.Transfer) (rpc.Transfer_Result, error) {
	// Initialize response object
	var response rpc.Transfer_Result

	// Define payload data
	payloadData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "transfer",
		"params": map[string]interface{}{
			"ringsize":  8, // Consider defining Ringsize as a constant
			"transfers": transfers,
			"scid":      exports.DERO_SCID_STRING,
		},
	}

	// Marshal payload data to JSON
	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		return response, fmt.Errorf("error marshaling payload JSON: %v", err)
	}

	// Define HTTP client and request
	client := &http.Client{}
	url := exports.DeroEndpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return response, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.SetBasicAuth(exports.DeroUsername, exports.DeroPassword)
	req.Header.Set("Content-Type", "application/json")

	// Send HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return response, fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read and process response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("error reading response body: %v", err)
	}

	// Print response body for debugging
	fmt.Printf("Response Body: %s\n", responseBody)

	// Unmarshal response JSON
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, fmt.Errorf("error unmarshaling response JSON: %v", err)
	}

	// Print TXID for reference
	fmt.Printf("TXID: %s", response.TXID)

	return response, nil
}
