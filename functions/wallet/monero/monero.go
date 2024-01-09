package monero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gabstv/httpdigest"
)

var (
	moneroUser = "secretnamebasis"
	moneroPass = "bargraph-chivalry-bullhorn"
	ip         = "192.168.12.176"
	port       = "28088"

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

// makeRequest sends an HTTP request and returns the response.
func makeRequest(request *http.Request) (*http.Response, error) {
	client := http.Client{
		Transport: &httpdigest.Transport{
			Username:  moneroUser,
			Password:  moneroPass,
			Transport: http.DefaultTransport,
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// createHTTPRequest creates an HTTP request with common headers and parameters.
func createHTTPRequest(method, endpoint string, body []byte) (*http.Request, error) {
	url := fmt.Sprintf("http://%s:%s/%s", ip, port, endpoint)
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
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
			"in": true,

			"filter_by_height": true,
			"min_height":       uint64(n - 1),
			"max_height":       uint64(n),
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
