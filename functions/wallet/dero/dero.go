package dero

import (
	"errors"
	"fmt"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
)

var err error
var clone *rpc.Address
var transfers rpc.Get_Transfers_Result

func Connection() bool {
	s := dero.Address()
	test := Echo(s)
	if test != "WALLET "+s+"\n" {
		return false
	}
	return true
}

func Height() int {
	err = exports.DeroRpcClient.CallFor(&exports.WalletHeight, "GetHeight")
	if err != nil || exports.WalletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(exports.WalletHeight.Height)
}

func Address() string {

	err = exports.DeroRpcClient.CallFor(&exports.Addr_result, "GetAddress")
	if err != nil || exports.Addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return err.Error()
	}

	exports.Addr, err = rpc.NewAddress(exports.Addr_result.Address)
	if err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", exports.Addr_result.Address, err)
		return err.Error()
	}
	return exports.Addr.String()
}

func SendTransfer(params rpc.Transfer_Params) string {
	var transfers rpc.Transfer_Result
	_ = exports.DeroRpcClient.CallFor(
		&transfers,
		"Transfer",
		params,
	)

	if transfers.TXID == "" {
		err := errors.New("Empty TXID")
		exports.Logs.Error(err, Echo("TXID is \"\" string"))
	}
	return transfers.TXID
}

func GetIncomingTransfers() (rpc.Get_Transfers_Result, error) {

	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In: true,
		},
	)
	if err != nil {
		exports.Logs.Error(err, "Could not obtain gettransfers from wallet")
		return transfers, err
	}

	return transfers, nil
}

func GetIncomingTransfersByHeight(h int) (*rpc.Get_Transfers_Result, error) {
	var transfers rpc.Get_Transfers_Result

	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:         true,
			Min_Height: uint64(h),
			Max_Height: uint64(h),
		},
	)
	if err != nil {
		exports.Logs.Error(err, "Could not obtain gettransfers from wallet")
		return nil, err
	}

	if len(transfers.Entries) == 0 {
		return nil, nil
	}

	return &transfers, nil
}

func CreateServiceAddress(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address := clone.Clone()
	service_address.Arguments = exports.Expected_arguments
	return service_address.String()
}

func CreateServiceAddressWithoutHardcodedValue(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address_without_amount := clone.Clone()

	service_address_without_amount.
		Arguments = exports.Expected_arguments[:len(exports.Expected_arguments)-1]

	return service_address_without_amount.String()
}

func Echo(s string) string {
	var echoResult string
	err := exports.DeroRpcClient.CallFor(&echoResult, "Echo", s+"\n")
	if err != nil {
		return err.Error()
	}

	return echoResult
}
