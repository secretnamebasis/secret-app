package functions

import (
	"fmt"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
)

var err error
var clone *rpc.Address

func Connection() bool {
	test := Echo(exports.Username)
	if test != "WALLET "+exports.Username+"\n" {
		return false
	}
	return true
}

func Height() int {
	err = exports.RpcClient.CallFor(&exports.WalletHeight, "GetHeight")
	if err != nil || exports.WalletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(exports.WalletHeight.Height)
}

func Address() string {

	err = exports.RpcClient.CallFor(&exports.Addr_result, "GetAddress")
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

func GetTransfers() (rpc.Get_Transfers_Result, error) {

	err = exports.RpcClient.CallFor(
		&exports.Transfers,
		"GetTransfers",
		rpc.Get_Transfers_Params{
			In: true,
		},
	)
	if err != nil {
		exports.Logs.Error(err, "Could not obtain gettransfers from wallet")
		return exports.Transfers, err
	}

	return exports.Transfers, nil
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
	err := exports.RpcClient.CallFor(&echoResult, "Echo", s+"\n")
	if err != nil {
		return err.Error()
	}

	return echoResult
}
