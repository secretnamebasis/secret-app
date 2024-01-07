package wallet

import (
	"strconv"
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/logger"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

var currentDeroHeight int
var currentMoneroHeight int

func IncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {

	checkAndProcess := func(transfers *rpc.Get_Transfers_Result) error {

		if !*LoopActivated {

			*LoopActivated = true
		}

		for _, e := range transfers.Entries {
			if err := dero.IncomingTransferEntry(e, db); err != nil {
				return err
			}
		}

		return nil
	}

	for {
		deroHeight := dero.Height()

		moneroHeight := monero.Height()

		if currentDeroHeight != deroHeight {
			currentDeroHeight = deroHeight

			if currentMoneroHeight != moneroHeight {
				currentMoneroHeight = moneroHeight
				logger.HeightInfo("Monero Height:" + strconv.Itoa(monero.Height()))

			}

			logger.HeightInfo("Dero Height:" + strconv.Itoa(dero.Height()))
			transfers, err := dero.GetIncomingTransfersByHeight(dero.Height())

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
