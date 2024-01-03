package exports

import (
	"net/http"

	"github.com/deroproject/derohe/rpc"
	"github.com/go-logr/logr"
	"github.com/secretnamebasis/secret-app/code"
	"github.com/ybbus/jsonrpc"
)

const DEVELOPER_ADDRESS = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
const DEVELOPER_NAME = "secretnamebasis"
const APP_NAME = "secret-app"
const DEST_PORT = uint64(0x1337)
const MESSAGE = "secret loves you"

var (
	Testing bool // Global variable to indicate testing mode
	err     error
	db_name string
	sale    []byte
)

var (
	username = "secret"
	password = "pass"
	ip       = "192.168.12.208"
	port     = "10104"
	endpoint = "http://" + ip + ":" + port + "/json_rpc"
)

var (
	Pong         = "You have purchased a really cool link"
	Welcome      = code.SayHello(username)
	walletHeight *rpc.GetHeight_Result
	addr         *rpc.Address
	clone        *rpc.Address
	addr_result  rpc.GetAddress_Result
	logger       logr.Logger = logr.Discard() // default discard all logs

	transfers rpc.Get_Transfers_Result
)

var (
	HttpClient = &http.Client{
		Transport: &TransportWithBasicAuth{
			Username: username,
			Password: password,
			Base:     http.DefaultTransport,
		},
	}
)

var (
	rpcClient = jsonrpc.NewClientWithOpts(
		endpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: HttpClient,
		},
	)
)

var (
	// currently the interpreter seems to have a glitch if this gets initialized within the code
	// see limitations github.com/traefik/yaegi
	response = rpc.Arguments{
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
	expected_arguments = rpc.Arguments{
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
