package handlers

import (
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/logger"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"go.etcd.io/bbolt"
)

func processIncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {

	checkAndProcess := func(transfers *rpc.Get_Transfers_Result) error {

		if !*LoopActivated {

			*LoopActivated = true
		}

		for _, e := range transfers.Entries {
			if err := IncomingTransferEntry(e, db); err != nil {
				return err
			}
		}

		return nil
	}

	for {
		height := dero.Height()

		if currentHeight != height {

			currentHeight = height
			logger.HeightInfo(strconv.Itoa(dero.Height()))

			transfers, err := dero.GetIncomingTransfersByHeight(currentHeight)

			if transfers == nil {
				continue
			}

			if err != nil {
				return err
			}

			if err := checkAndProcess(transfers); err != nil {
				return err
			}
		}

		sleepDuration := 1 * time.Second
		if exports.Testing {
			sleepDuration = 100 * time.Millisecond
		}
		time.Sleep(sleepDuration)
	}
}
