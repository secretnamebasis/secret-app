package generators

import (
	"github.com/deroproject/derohe/rpc"
)

func Transfers(destinations []string, amounts []uint64, responses []rpc.Arguments) []rpc.Transfer {
	if len(destinations) != len(amounts) || len(destinations) != len(responses) {
		// Return an empty slice of transfers if lengths of provided slices are not equal
		return []rpc.Transfer{}
	}

	var transfers []rpc.Transfer

	for i, dest := range destinations {
		transfers = append(transfers, rpc.Transfer{
			Destination: dest,
			Amount:      amounts[i],
			Payload_RPC: responses[i],
		})
	}
	return transfers
}
