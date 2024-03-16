package dero

import (
	"fmt"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
)

var err error
var clone *rpc.Address
var transfers rpc.Get_Transfers_Result
var balance rpc.GetBalance_Result

func Connection() (bool, error) {

	_, err := Address()

	if err != nil {
		return false, err
	}
	return true, nil
}

func Height() int {
	err = exports.DeroRpcClient.CallFor(
		&exports.WalletHeight,
		"GetHeight",
	)
	if err != nil || exports.WalletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(exports.WalletHeight.Height)
}

func Address() (string, error) {

	err = exports.DeroRpcClient.CallFor(&exports.Addr_result, "GetAddress")
	if err != nil || exports.Addr_result.Address == "" {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return err.Error(), err
	}

	exports.Addr, err = rpc.NewAddress(exports.Addr_result.Address)
	if err != nil {
		fmt.Printf("address could not be parsed: addr:%s err:%s\n", exports.Addr_result.Address, err)
		return err.Error(), err
	}
	return exports.Addr.String(), nil
}

func SendTransfer(params rpc.Transfer_Params) (rpc.Transfer_Result, error) {
	var transfers rpc.Transfer_Result
	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"Transfer",
		params,
	)

	if err != nil {
		return transfers, err
	}

	return transfers, nil
}

func GetBalance() (rpc.GetBalance_Result, error) {

	err = exports.DeroRpcClient.CallFor(
		&balance,
		"GetBalance",
	)
	return balance, nil
}

func GetIncomingTransfers() (rpc.Get_Transfers_Result, error) {

	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:       true,
			Coinbase: false,
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

func GetOutgoingTransfers() (rpc.Get_Transfers_Result, error) {

	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:       false,
			Out:      true,
			Coinbase: false,
		},
	)
	if err != nil {
		exports.Logs.Error(err, "Could not obtain gettransfers from wallet")
		return transfers, err
	}

	return transfers, nil
}

func GetAllTransfers() (rpc.Get_Transfers_Result, error) {

	err = exports.DeroRpcClient.CallFor(
		&transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In:       true,
			Out:      true,
			Coinbase: false,
		},
	)
	if err != nil {
		exports.Logs.Error(err, "Could not obtain gettransfers from wallet")
		return transfers, err
	}

	return transfers, nil
}

func CreateServiceAddress(order string) (string, error) {
	s, err := Address()
	if err != nil {
		return "", err
	}
	clone, err = rpc.NewAddress(s)
	service_address := clone.Clone()
	exports.Expected_arguments[1].Value = order
	service_address.Arguments = exports.Expected_arguments
	return service_address.String(), nil
}

func CreateServiceAddressWithoutHardcodedValue(addr string) (string, error) {
	s, err := Address()
	if err != nil {
		return "", err
	}
	clone, err = rpc.NewAddress(s)
	service_address_without_amount := clone.Clone()

	service_address_without_amount.
		Arguments = exports.Expected_arguments[:len(exports.Expected_arguments)-1]

	return service_address_without_amount.String(), nil
}

func Echo(s string) string {
	var echoResult string
	err := exports.DeroRpcClient.CallFor(&echoResult, "Echo", s+"\n")
	if err != nil {
		return err.Error()
	}

	return echoResult
}
