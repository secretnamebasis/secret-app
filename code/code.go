package code

import (
	"crypto/sha1"
	"errors"
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
	addr         *rpc.Address
	clone        *rpc.Address
	addr_result  rpc.GetAddress_Result
	logger       logr.Logger = logr.Discard() // default discard all logs
	sale                     = []byte("")
	username                 = Name
	password                 = "pass"
	endpoint                 = "http://192.168.12.208:10104/json_rpc"
	db_name      string

	transfers  rpc.Get_Transfers_Result
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

	if WalletConnection() == false {
		return fmt.Errorf("Failed to establish wallet connection")
	}

	Logger()

	fmt.Printf(WalletEcho("Logger has started"))

	// Let's make a database
	db_name = fmt.Sprintf("%s_%s.bbolt.db", APP_NAME, Sha1Sum(WalletAddress()))
	fmt.Printf(WalletEcho("ID has been created"))

	go func() { // flip that shit on
		db, err := CreateDB(db_name)
		defer db.Close()
		if err != nil {
			logger.Error(err, err.Error())
		}
		fmt.Printf(WalletEcho("Database has been created"))

		// Let's make a bucket
		sale = []byte("SALE")
		CreateBucket(db, sale)
		fmt.Printf(WalletEcho("Sale's list initiated"))
	}()

	transfers, err := WalletGetTransfers()
	if err != nil {
		// Print error and return
		logger.Error(err, "Error getting transfers from wallet")
		return err
	}
	fmt.Printf(WalletEcho("Transfers retreived"))
	_, err = HandleIncomingTransfers(transfers.Entries)
	if err != nil {

		// Print error and return
		logger.Error(err, "Error handling incoming transfers")
		return err
	}
	// fmt.Printf(handle)

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

	db, err := bbolt.Open(db_name, 0600, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	return db, nil
}

func WalletConnection() bool {
	test := WalletEcho(Name)
	if test != "WALLET "+Name+"\n" {
		return false
	}
	return true
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

	err = RpcClient.CallFor(&addr_result, "GetAddress")
	if err != nil || addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return err.Error()
	}

	addr, err = rpc.NewAddress(addr_result.Address)
	if err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", addr_result.Address, err)
		return err.Error()
	}
	return addr.String()
}

func WalletGetTransfers() (rpc.Get_Transfers_Result, error) {

	err = RpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In: true,
		},
	)
	if err != nil {
		logger.Error(err, "Could not obtain gettransfers from wallet")
		return transfers, err
	}

	return transfers, nil
}

func HandleIncomingTransfers(entries []rpc.Entry) (string, error) {
	var errorMsg string

	for _, e := range entries {
		// Simulating a condition that might lead to an error
		if e.Amount <= 0 {
			return fmt.Sprintf("invalid transaction amount: %d", e.Amount), errors.New("invalid transaction amount")
		}

		// Preparing the processing message
		msg := fmt.Sprintf("Processing incoming transaction: TXID - %s, Amount - %d\n", e.TXID, e.Amount)
		errorMsg += msg
	}

	return errorMsg, nil
}

func CreateServiceAddress(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address := clone.Clone()
	return service_address.String()
}

func CreateServiceAddressWithoutHardcodedValue(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address_without_amount := clone.Clone()

	service_address_without_amount.
		Arguments = Expected_arguments[:len(Expected_arguments)-1]

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
