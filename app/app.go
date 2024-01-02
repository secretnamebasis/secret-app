package app

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"

	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/go-logr/logr"
	"github.com/ybbus/jsonrpc"
	"go.etcd.io/bbolt"
	"gopkg.in/natefinch/lumberjack.v2"
)

const DEVELOPER_ADDRESS = "dero1qyvqpdftj8r6005xs20rnflakmwa5pdxg9vcjzdcuywq2t8skqhvwqglt6x0g"
const DEVELOPER_NAME = "secretnamebasis"
const APP_NAME = "secret-app"
const DEST_PORT = uint64(0x1337)
const MESSAGE = "secret loves you"

var Testing = false // Global variable to indicate testing mode
var (
	err          error
	Name         = "secret"
	Welcome      = SayHello(Name)
	walletHeight *rpc.GetHeight_Result
	Addr         *rpc.Address
	Addr_result  rpc.GetAddress_Result
	logger       logr.Logger = logr.Discard() // default discard all logs
	db_name                  = fmt.Sprintf("%s_%s.bbolt.db", APP_NAME, Sha1Sum(WalletAddress()))
	sale                     = []byte("SALE")
	username                 = Name
	password                 = "pass"
	endpoint                 = "http://192.168.12.208:10104/json_rpc"

	HttpClient = &http.Client{
		Transport: &TransportWithBasicAuth{
			Username: username,
			Password: password,
			Base:     http.DefaultTransport,
		},
	}
	RpcClient = jsonrpc.NewClientWithOpts(
		endpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient: HttpClient,
		},
	)

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

func RunApp() error {
	Logger()
	fmt.Printf(WalletEcho(Name))
	db, err := CreateDB(db_name)
	if err != nil {
		logger.Error(err, err.Error())
	}
	CreateBucket(db, sale)
	CreateServiceAddress()
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

func Sha1Sum(DEVELOPER_ADDRESS string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(WalletAddress())))
	return shasum
}

func Logger() error {
	// parse arguments and setup logging, print basic information
	globals.Arguments["--debug"] = true
	exename, err := os.Executable()
	if err != nil {
		globals.InitializeLog(os.Stdout, &lumberjack.Logger{
			Filename:   exename + ".log",
			MaxSize:    100, // megabytes
			MaxBackups: 2,
		})
		logger = globals.Logger
		logger.Error(err, "Logger failed")
		return err
	}

	return nil
}

func CreateBucket(db *bbolt.DB, bucketName []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func CreateDB(db_name string) (*bbolt.DB, error) {
	if db_name == "" {
		db_name = fmt.Sprintf("%s_%s.bbolt.db", APP_NAME, Sha1Sum(WalletAddress()))
	}
	if !Testing {
		db, err := bbolt.Open(db_name, 0600, nil)
		if err != nil {
			fmt.Printf("could not open db err:%s\n", err)
			return nil, err
		}

		return db, nil
	}

	return nil, nil
}

func WalletHeight() int {
	err = RpcClient.CallFor(&walletHeight, "GetHeight")
	if err != nil || walletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(walletHeight.Height)
}

func WalletAddress() string {

	err = RpcClient.CallFor(&Addr_result, "GetAddress")
	if err != nil || Addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return err.Error()
	}

	if Addr, err = rpc.NewAddress(Addr_result.Address); err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", Addr_result.Address, err)
		return err.Error()
	}
	return Addr.String()
}
func CreateServiceAddress() string {
	service_address_without_amount := Addr.Clone()
	service_address_without_amount.Arguments = Expected_arguments[:len(Expected_arguments)-1]
	err, _ := fmt.Printf("Integrated address to activate '%s'\nWithout hardcoded amount) service: \n%s\n", APP_NAME, service_address_without_amount.String())
	if err == 0 {
		return ""
	}
	return service_address_without_amount.String()
}

func WalletEcho(Name string) string {
	var echoResult string
	err := RpcClient.CallFor(&echoResult, "Echo", Name+"\n")
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
