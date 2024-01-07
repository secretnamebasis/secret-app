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

func MakeIntegratedAddress() (map[string]string, error) {
	params := map[string]interface{}{}

	data := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "0",
		"method":  "make_integrated_address",
		"params":  params,
	}
	jsonData, _ := json.Marshal(data)
	url := fmt.Sprintf("http://%s:%s/json_rpc", ip, port)

	client := http.Client{
		Transport: &httpdigest.Transport{
			Username:  moneroUser,
			Password:  moneroPass,
			Transport: http.DefaultTransport,
		},
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseJSON IntegratedAddressResponse
	err = json.NewDecoder(response.Body).Decode(&responseJSON)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	result["integrated_address"] = responseJSON.Result.IntegratedAddress
	result["payment_id"] = responseJSON.Result.PaymentID

	return result, nil
}
