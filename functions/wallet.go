package functions

import (
	"fmt"

	"github.com/deroproject/derohe/rpc"
)

func Connection() bool {
	test := Echo(username)
	if test != "WALLET "+username+"\n" {
		return false
	}
	return true
}

func WalletHeight() int {
	err = rpcClient.CallFor(&walletHeight, "GetHeight")
	if err != nil || walletHeight.Height == 0 {
		fmt.Printf("Could not obtain address from wallet err %s\n", err)
		return 0
	}
	return int(walletHeight.Height)
}

func Address() string {

	err = rpcClient.CallFor(&addr_result, "GetAddress")
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

func GetTransfers() (rpc.Get_Transfers_Result, error) {

	err = rpcClient.CallFor(
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

func CreateServiceAddress(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address := clone.Clone()
	service_address.Arguments = expected_arguments
	return service_address.String()
}

func CreateServiceAddressWithoutHardcodedValue(addr string) string {
	clone, err = rpc.NewAddress(addr)
	service_address_without_amount := clone.Clone()

	service_address_without_amount.
		Arguments = expected_arguments[:len(expected_arguments)-1]

	return service_address_without_amount.String()
}

func Echo(username string) string {
	var echoResult string
	err := rpcClient.CallFor(&echoResult, "Echo", username+"\n")
	if err != nil {
		return err.Error()
	}

	return echoResult
}
