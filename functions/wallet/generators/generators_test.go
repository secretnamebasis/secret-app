package generators_test

import (
	"testing"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/generators"
)

func TestGenerators(t *testing.T) {
	t.Run(
		"TestTransfers",
		func(t *testing.T) {
			where := []string{
				exports.SECRET_WALLET_ADDRESS,
				exports.CAPTAIN_WALLET_ADDRESS,
				exports.B_WALDO_WALLET_ADDRESS,
			}
			count := []uint64{
				1,
				1,
				1,
			}
			response := []rpc.Arguments{
				exports.Response,
				exports.Response,
				exports.Response,
			}

			got := generators.Transfers(where, count, response)

			var result rpc.Transfer_Result

			err := exports.DeroRpcClient.CallFor(
				&result,
				"Transfer",
				rpc.Transfer_Params{Transfers: got},
			)

			if err != nil {
				t.Errorf("error while calling Transfer: %v", err)
			}

			if result.TXID == "" {
				t.Error("expected non-empty TXID, got empty")
			}
		},
	)

}
