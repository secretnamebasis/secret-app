package exports

import (
	"net/http"

	"github.com/deroproject/derohe/rpc"

	"github.com/ybbus/jsonrpc"
)

var (
	Testing bool // Global variable to indicate testing mode

)

var (
	Username = "secret"
	password = "pass"
	ip       = "192.168.12.208"
	port     = "10104"
	endpoint = "http://" + ip + ":" + port + "/json_rpc"
)

var (
	Pong         = "You have purchased a really cool link"
	Welcome      string
	WalletHeight *rpc.GetHeight_Result
	Addr         *rpc.Address

	Addr_result rpc.GetAddress_Result
	// default discard all logs

	Transfers rpc.Get_Transfers_Result
)

var (
	HttpClient = &http.Client{
		Transport: &TransportWithBasicAuth{
			Username: Username,
			Password: password,
			Base:     http.DefaultTransport,
		},
	}
)

var (
	RpcClient = jsonrpc.NewClientWithOpts(
		endpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: HttpClient,
		},
	)
)

var (
	// currently the interpreter seems to have a glitch if this gets initialized within the code
	// see limitations github.com/traefik/yaegi
	Response = rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		},
		{
			Name:     rpc.RPC_SOURCE_PORT,
			DataType: rpc.DataUint64,
			Value:    DEST_PORT,
		},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    "Successfully purchased pong (this could be serial/license key or download link or further)",
		},
	}
)

var (
	Expected_arguments = rpc.Arguments{
		{
			Name:     rpc.RPC_DESTINATION_PORT,
			DataType: rpc.DataUint64,
			Value:    DEST_PORT,
		},
		// { Name:rpc.RPC_EXPIRY , DataType:rpc.DataTime, Value:time.Now().Add(time.Hour).UTC()},
		{
			Name:     rpc.RPC_COMMENT,
			DataType: rpc.DataString,
			Value:    MESSAGE,
		},
		//{"float64", rpc.DataFloat64, float64(0.12345)},          // in atomic units
		{
			Name:     rpc.RPC_NEEDS_REPLYBACK_ADDRESS,
			DataType: rpc.DataUint64,
			Value:    uint64(0),
		}, // this service will reply to incoming request,so needs the senders address
		{
			Name:     rpc.RPC_VALUE_TRANSFER,
			DataType: rpc.DataUint64,
			Value:    uint64(12345),
		}, // in atomic units

	}
)
