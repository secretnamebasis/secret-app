package exports

import (
	"log"
	"net/http"
	"os"

	"github.com/deroproject/derohe/rpc"
	"github.com/gabstv/httpdigest"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/ybbus/jsonrpc"
)

var (
	Testing bool // Global variable to indicate testing mode

	Pong         = "You have purchased a really cool link"
	deroUsername string
	deroPassword string
	deroIp       string
	deroPort     string
	deroEndpoint string

	moneroUsername string
	moneroPassword string
	moneroIp       string
	moneroPort     string
	MoneroEndpoint string

	Welcome      string
	WalletHeight *rpc.GetHeight_Result
	Addr         *rpc.Address

	Addr_result rpc.GetAddress_Result

	DeroHttpClient *http.Client

	DeroRpcClient jsonrpc.RPCClient

	MoneroHttpClient *http.Client

	MoneroRpcClient jsonrpc.RPCClient

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
			Value:    uuid.New().String(),
		},
	}

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
			Value:    PONG_AMOUNT,
		}, // in atomic units

	}
)

func init() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	deroUsername = os.Getenv("DERO_WALLET_USER")
	deroPassword = os.Getenv("DERO_WALLET_PASS")
	deroIp = os.Getenv("DERO_SERVER_IP")
	deroPort = os.Getenv("DERO_WALLET_PORT")
	deroEndpoint = "http://" + deroIp + ":" + deroPort + "/json_rpc"
	moneroUsername = os.Getenv("MONERO_WALLET_USER")
	moneroPassword = os.Getenv("MONERO_WALLET_PASS")
	moneroIp = os.Getenv("MONERO_SERVER_IP")
	moneroPort = os.Getenv("MONERO_WALLET_PORT")
	MoneroEndpoint = "http://" + moneroIp + ":" + moneroPort + "/json_rpc"

	DeroHttpClient = &http.Client{
		Transport: &TransportWithBasicAuth{
			Username: deroUsername,
			Password: deroPassword,
			Base:     http.DefaultTransport,
		},
	}
	DeroRpcClient = jsonrpc.NewClientWithOpts(
		deroEndpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: DeroHttpClient,
		},
	)

	MoneroHttpClient = &http.Client{
		Transport: &httpdigest.Transport{
			Username:  moneroUsername,
			Password:  moneroPassword,
			Transport: http.DefaultTransport,
		},
	}

	MoneroRpcClient = jsonrpc.NewClientWithOpts(
		MoneroEndpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: MoneroHttpClient,
		},
	)
}
