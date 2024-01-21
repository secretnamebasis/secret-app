package monero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/secretnamebasis/secret-app/exports"
)

var (
	in = true
)

type (
	TransferResult struct {
		In      []Entry `json:"in,omitempty"`
		Out     []Entry `json:"out,omitempty"`
		Pending []Entry `json:"pending,omitempty"`
		Failed  []Entry `json:"failed,omitempty"`
		Pool    []Entry `json:"pool,omitempty"`
	}

	SubaddrIndex struct {
		Major uint64 `json:"major"`
		Minor uint64 `json:"minor"`
	}
	Entry struct {
		Address                         string         `json:"address"`
		Amount                          uint64         `json:"amount"`
		Amounts                         []uint64       `json:"amounts"`
		Confirmations                   uint64         `json:"confirmations"`
		DoubleSpendSeen                 bool           `json:"double_spend_seen"`
		Fee                             uint64         `json:"fee"`
		Height                          uint64         `json:"height"`
		Locked                          bool           `json:"locked"`
		Note                            string         `json:"note"`
		PaymentID                       string         `json:"payment_id"`
		SubaddrIndex                    SubaddrIndex   `json:"subaddr_index"`
		SubaddrIndices                  []SubaddrIndex `json:"subaddr_indices"`
		SuggestedConfirmationsThreshold uint64         `json:"suggested_confirmations_threshold"`
		Timestamp                       uint64         `json:"timestamp"`
		TxID                            string         `json:"txid"`
		Type                            string         `json:"type"`
		UnlockTime                      uint64         `json:"unlock_time"`
	}

	Get_Transfers_Result struct {
		Entries TransferResult
	}
	IntegratedAddressResponse struct {
		ID      string `json:"id"`
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			IntegratedAddress string `json:"integrated_address"`
			PaymentID         string `json:"payment_id"`
		} `json:"result"`

		// Add the fields IntegratedAddress and PaymentID directly in the struct
		IntegratedAddress string
		PaymentID         string
	}
	HeightResponse struct {
		ID      string `json:"id"`
		JSONRPC string `json:"jsonrpc"`
		Result  struct {
			Height uint64 `json:"height"`
		} `json:"result"`
	}

	Transfer_Params struct {
		Destinations   []Transfer `json:"destinations"`
		AccountIndex   uint       `json:"account_index,omitempty"`
		SubaddrIndices []uint     `json:"subaddr_indices,omitempty"`
		Priority       uint       `json:"priority,omitempty"`
		Mixin          uint       `json:"mixin,omitempty"`
		RingSize       uint       `json:"ring_size,omitempty"`
		UnlockTime     uint       `json:"unlock_time,omitempty"`
		GetTxKey       bool       `json:"get_tx_key,omitempty"`
		DoNotRelay     bool       `json:"do_not_relay,omitempty"`
		GetTxHex       bool       `json:"get_tx_hex,omitempty"`
		GetTxMetadata  bool       `json:"get_tx_metadata,omitempty"`
	}

	Transfer struct {
		Amount         uint64 `json:"amount"`
		Address        string `json:"address"`
		AccountIndex   uint   `json:"account_index,omitempty"`
		SubaddrIndices []uint `json:"subaddr_indices,omitempty"`
		Priority       uint   `json:"priority,omitempty"`
		Mixin          uint   `json:"mixin,omitempty"`
		RingSize       uint   `json:"ring_size,omitempty"`
		UnlockTime     uint   `json:"unlock_time,omitempty"`
		GetTxKey       bool   `json:"get_tx_key,omitempty"`
		DoNotRelay     bool   `json:"do_not_relay,omitempty"`
		GetTxHex       bool   `json:"get_tx_hex,omitempty"`
		GetTxMetadata  bool   `json:"get_tx_metadata,omitempty"`
	}

	TransferResponse struct {
		Result Result `json:"result"`
	}

	Result struct {
		Amount        uint64 `json:"amount"`
		Fee           uint64 `json:"fee"`
		MultisigTxset string `json:"multisig_txset"`
		TxBlob        string `json:"tx_blob"`
		TxHash        string `json:"tx_hash"`
		TxKey         string `json:"tx_key"`
		TxMetadata    string `json:"tx_metadata"`
		UnsignedTxset string `json:"unsigned_txset"`
	}
)

func Height() int {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_height",
	}
	jsonData, _ := json.Marshal(data)

	request, _ := createHTTPRequest("POST", "json_rpc", jsonData)
	if request == nil {
		return 0
	}
	defer request.Body.Close()

	response, _ := makeRequest(request)
	if response == nil {
		return 0
	}
	defer response.Body.Close()

	var responseJSON HeightResponse
	err := json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		return 0
	}

	return int(responseJSON.Result.Height)
}

// MakeIntegratedAddress generates an integrated address.
func MakeIntegratedAddress() (map[string]string, error) {
	params := map[string]interface{}{}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "make_integrated_address",
		"params":  params,
	}
	jsonData, _ := json.Marshal(data)

	request, err := createHTTPRequest("POST", "json_rpc", jsonData)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	response, err := makeRequest(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseJSON IntegratedAddressResponse
	if err := json.NewDecoder(response.Body).Decode(&responseJSON); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	result["integrated_address"] = responseJSON.Result.IntegratedAddress
	result["payment_id"] = responseJSON.Result.PaymentID

	return result, nil
}

func GetIncomingTransfers() (TransferResult, error) {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_transfers",
		"params": map[string]interface{}{
			"in": in,
		},
	}
	jsonData, _ := json.Marshal(data)

	request, err := createHTTPRequest("POST", "json_rpc", jsonData)
	if err != nil {
		return TransferResult{}, err
	}
	defer request.Body.Close()

	response, err := makeRequest(request)
	if err != nil {
		return TransferResult{}, err
	}
	defer response.Body.Close()

	var transferResponse struct {
		Result TransferResult `json:"result,omitempty"`
	}
	if err := json.NewDecoder(response.Body).Decode(&transferResponse); err != nil {
		return TransferResult{}, err
	}

	return transferResponse.Result, nil
}

func GetIncomingTransfersByHeight(n int) (TransferResult, error) {
	/*
		NOTE: this doesn't work as intended
		For what ever reason, and ignorance
		is a valid reason for reason...
		When the server is called, it
		returns the appropriate height object
		from result:{"in":[]} (see tests);
		however, despite my best efforts,
		the transfers at Height(),
		will return nil always.

		So our solution was to dump all
		transfers through the handler and
		reference already_processed for
		all objects in the wallet that
		have... already been processed.
	*/
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_transfers",
		"params": map[string]interface{}{
			"in":               true,
			"filter_by_height": true,
			"min_height":       n - 1000,
			"max_height":       n,
		},
	}
	jsonData, _ := json.Marshal(data)

	request, err := createHTTPRequest("POST", "json_rpc", jsonData)
	if err != nil {
		return TransferResult{}, err
	}
	defer request.Body.Close()

	response, err := makeRequest(request)
	if err != nil {
		return TransferResult{}, err
	}
	defer response.Body.Close()

	// Create a new io.Reader for the response body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return TransferResult{}, err
	}

	// Create a new io.Reader using the bodyBytes
	bodyReader := bytes.NewReader(bodyBytes)

	var transferResponse struct {
		Result TransferResult `json:"result,omitempty"`
	}
	if err := json.NewDecoder(bodyReader).Decode(&transferResponse); err != nil {
		return TransferResult{}, err
	}

	return transferResponse.Result, nil
}

// Address retrieves the wallet's addresses for a specific account index.
func Address(accountIndex uint64) string {
	params := map[string]interface{}{

		// Optionally, add logic to specify specific address indices if needed.
		// "address_index": []uint64{0, 1, 4},
	}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "get_address",
		"params":  params,
	}
	jsonData, _ := json.Marshal(data)

	request, _ := createHTTPRequest("POST", "json_rpc", jsonData)

	defer request.Body.Close()

	response, _ := makeRequest(request)

	defer response.Body.Close()

	var addressResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&addressResponse); err != nil {
		return err.Error()
	}

	// Extract the primary address from the response
	address := addressResponse["result"].(map[string]interface{})["address"].(string)
	return address
}

func ValidateAddress(address string) bool {
	params := map[string]interface{}{
		"address": address,
	}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "validate_address",
		"params":  params,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	// Assuming createHTTPRequest, makeRequest, and response handling functions are defined elsewhere
	request, _ := createHTTPRequest("POST", "json_rpc", jsonData)
	defer request.Body.Close()

	response, _ := makeRequest(request)
	defer response.Body.Close()

	var validateResponse map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&validateResponse); err != nil {
		return false
	}

	// Extract the 'valid' field from the response
	isValid := validateResponse["result"].(map[string]interface{})["valid"].(bool)
	return isValid
}

func SendTransfer(params Transfer_Params) TransferResponse {
	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "transfer",
		"params":  params,
	}
	jsonData, _ := json.Marshal(data)

	request, _ := createHTTPRequest("POST", "json_rpc", jsonData)
	defer request.Body.Close()

	response, _ := makeRequest(request)
	defer response.Body.Close()

	// Create a new io.Reader for the response body
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return TransferResponse{}
	}

	// Print or log the entire response body for debugging purposes
	fmt.Printf("Response Body: %s\n", bodyBytes)

	// Create a new io.Reader using the bodyBytes
	bodyReader := bytes.NewReader(bodyBytes)

	var transferResponse TransferResponse
	if err := json.NewDecoder(bodyReader).Decode(&transferResponse); err != nil {
		fmt.Printf("Error decoding JSON response: %v\n", err)
		return TransferResponse{}
	}

	return transferResponse
}

// private methods

// makeRequest sends an HTTP request and returns the response.
func makeRequest(request *http.Request) (*http.Response, error) {
	client := exports.MoneroHttpClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// createHTTPRequest creates an HTTP request with common headers and parameters.
func createHTTPRequest(method, endpoint string, body []byte) (*http.Request, error) {
	request, err := http.NewRequest(method, exports.MoneroEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set common headers
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}
