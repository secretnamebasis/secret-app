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
)

// IntegratedAddressResponse represents the structure of the JSON response for an integrated address
type IntegratedAddressResponse struct {
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

type HeightResponse struct {
	ID      string `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		Height uint64 `json:"height"`
	} `json:"result"`
}

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
