package app

import (
	"fmt"
	"net/http"

	"github.com/ybbus/jsonrpc"
)

var Name = "secret"
var err error

func RunApp() error {
	fmt.Println(WalletEcho(Name))
	return nil
}

func Echo(Name string) string {
	return Name
}

func SayHello(Name string) string {
	return "Hello, " + Echo(Name)
}

func Ping() bool {
	return true
}

func WalletEcho(Name string) string {
	// Define basic auth username and password
	username := Name
	password := "pass"

	// Create a custom HTTP client with basic authentication
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Set basic authentication credentials
	httpClient.Transport = &TransportWithBasicAuth{
		Username: username,
		Password: password,
		Base:     httpClient.Transport,
	}

	// Create JSON-RPC client with custom HTTP client
	rpcClient := jsonrpc.NewClientWithOpts(
		"http://192.168.12.208:10104/json_rpc",
		&jsonrpc.RPCClientOpts{
			HTTPClient: httpClient,
		},
	)

	var echoResult string
	err := rpcClient.CallFor(&echoResult, "Echo", Name)
	if err != nil {
		return err.Error()
	}

	return echoResult
}

// TransportWithBasicAuth is a custom HTTP transport to set basic authentication
type TransportWithBasicAuth struct {
	Username string
	Password string
	Base     http.RoundTripper
}

// RoundTrip executes a single HTTP transaction, adding basic auth headers
func (t *TransportWithBasicAuth) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return t.Base.RoundTrip(req)
}
