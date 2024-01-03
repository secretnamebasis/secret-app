package code

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/go-logr/logr"
	"github.com/secretnamebasis/secret-app/code/exports"
	"github.com/secretnamebasis/secret-app/code/functions"
	"go.etcd.io/bbolt"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	db_name string
	sale    []byte
	logger  logr.Logger = logr.Discard()
)

func RunApp() error {

	Logger()

	if functions.Connection() == false {
		err := errors.New("Wallet Connection Failure")
		logger.Error(err, "Error")
		return fmt.Errorf(
			"Failed to establish wallet connection",
		)
	}
	logger.Info(
		functions.Echo(
			"Logger has started",
		),
	)

	// Let's make a database
	db_name = fmt.Sprintf(
		"%s_%s.bbolt.db",
		exports.APP_NAME,
		Sha1Sum(functions.Address()),
	)

	logger.Info(
		functions.Echo(
			"ID has been created",
		),
	)

	db, err := CreateDB(db_name)

	if err != nil {
		logger.Error(err, err.Error())
	}

	logger.Info(
		functions.Echo(
			"Database has been created",
		),
	)

	// Let's make a bucket
	sale = []byte("SALE")
	CreateBucket(db, sale)

	logger.Info(
		functions.Echo(
			"Sale's list initiated",
		),
	)
	logger.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments: " +
				CreateServiceAddress(
					functions.Address(),
				),
		),
	)

	logger.Info(
		functions.Echo(
			"Integrated Address with Expected Arguments minus Hardcoded Value: " +
				CreateServiceAddressWithoutHardcodedValue(
					functions.Address(),
				),
		),
	)

	HandleIncomingTransfers(db)

	return nil // Stop the loop and return nil

}

func Echo(username string) string {
	return username
}

func SayHello(username string) string {
	return "Hello, " + Echo(username)
}

func Ping() bool {
	return true
}

func Sha1Sum(DEVELOPER_ADDRESS string) string {
	shasum := fmt.Sprintf("%x", sha1.Sum([]byte(functions.Address())))
	return shasum
}

func Logger() error {
	// parse arguments and setup logging, print basic information
	globals.Arguments["--debug"] = true
	exename, err := os.Executable()

	globals.InitializeLog(os.Stdout, &lumberjack.Logger{
		Filename:   exename + ".log",
		MaxSize:    100, // megabytes
		MaxBackups: 2,
	})
	logger = globals.Logger

	return err

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

func HandleIncomingTransfers(db *bbolt.DB) error {
	forLoop := false
	for {
		transfers, err := functions.GetTransfers()
		if err != nil {
			logger.Error(err, "Wallet Failed to Get Entries")
		}
		if !forLoop {
			logger.Info("Wallet Entries are Instantiated")
			forLoop = true
		}

		for _, e := range transfers.Entries {
			// Simulating a condition that might lead to an error
			if e.Amount <= 0 {
				return errors.New("invalid transaction amount")
			}
			if e.Coinbase || !e.Incoming { // skip coinbase or outgoing, self generated transactions
				continue
			}

			// check whether the entry has been processed before, if yes skip it
			var already_processed bool
			db.View(func(tx *bbolt.Tx) error {
				if b := tx.Bucket([]byte("SALE")); b != nil {
					if ok := b.Get([]byte(e.TXID)); ok != nil { // if existing in bucket
						already_processed = true
					}
				}
				return nil
			})

			if already_processed { // if already processed skip it
				continue
			}

			// check whether this service should handle the transfer
			if !e.Payload_RPC.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) ||
				DEST_PORT != e.Payload_RPC.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64) { // this service is expecting value to be specfic
				continue

			}

			logger.V(1).Info("to be processed", "txid", e.TXID)
			if expected_arguments.Has(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64) { // this service is expecting value to be specfic
				value_expected := expected_arguments.Value(rpc.RPC_VALUE_TRANSFER, rpc.DataUint64).(uint64)
				if e.Amount != value_expected { // TODO we should mark it as faulty
					logger.Error(nil, fmt.Sprintf("user transferred %d, we were expecting %d. so we will not do anything", e.Amount, value_expected)) // this is an unexpected situation
					continue
				}

				if !e.Payload_RPC.Has(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress) {
					logger.Error(nil, fmt.Sprintf("user has not give his address so we cannot replyback")) // this is an unexpected situation
					continue
				}

				destination_expected := e.Payload_RPC.Value(rpc.RPC_REPLYBACK_ADDRESS, rpc.DataAddress).(rpc.Address).String()
				addr, err := rpc.NewAddress(destination_expected)
				if err != nil {
					logger.Error(err, "err while while parsing incoming addr")
					continue
				}
				addr.Mainnet = true // convert addresses to testnet form, by default it's expected to be mainnnet
				destination_expected = addr.String()

				logger.V(1).Info("tx should be replied", "txid", e.TXID, "replyback_address", destination_expected)

				//destination_expected := e.Sender

				// value received is what we are expecting, so time for response
				response[0].Value = e.SourcePort // source port now becomes destination port, similar to TCP
				response[2].Value = fmt.Sprintf("%s. You sent %s at height %d", Pong, walletapi.FormatMoney(e.Amount), e.Height)

				//_, err :=  response.CheckPack(transaction.PAYLOAD0_LIMIT)) //  we only have 144 bytes for RPC

				// sender of ping now becomes destination
				var result rpc.Transfer_Result
				tparams := rpc.Transfer_Params{Transfers: []rpc.Transfer{{Destination: destination_expected, Amount: uint64(1), Payload_RPC: response}}}
				err = rpcClient.CallFor(&result, "Transfer", tparams)
				if err != nil {
					logger.Error(err, "err while transfer")
					continue
				}

				err = db.Update(func(tx *bbolt.Tx) error {
					b := tx.Bucket([]byte("SALE"))
					return b.Put([]byte(e.TXID), []byte("done"))
				})
				if err != nil {
					logger.Error(err, "err updating db")
				} else {
					logger.Info("ping replied successfully with pong ", "result", result)
				}
				if Testing == true {
					return nil
				}
			}
		}

	}
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
