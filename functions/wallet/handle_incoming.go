package wallet

import (
	"time"

	"github.com/deroproject/derohe/rpc"
	"github.com/secretnamebasis/secret-app/exports"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"github.com/secretnamebasis/secret-app/functions/wallet/monero"
	"go.etcd.io/bbolt"
)

var currentDeroHeight int
var currentMoneroHeight int

func IncomingTransfers(db *bbolt.DB) error {
	LoopActivated := false

	if err := dero.Load(db); err != nil {
		return err // Exit on error during initial load
	}
	if err := monero.Load(db); err != nil {
		return err // Exit on error during initial load
	}
	Info("Entries Loaded")

	return processIncomingTransfers(db, &LoopActivated)
}

func processIncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {
	processMoneroTransfers := func(transfers *monero.Get_Transfers_Result) error {

		if !*LoopActivated {

			*LoopActivated = true
		}

		for _, e := range transfers.Entries.In {
			if err := monero.IncomingTransferEntry(e, db); err != nil {
				return err
			}
		}

		return nil
	}
	processDeroTransfers := func(transfers *rpc.Get_Transfers_Result) error {

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
	Info("Entering For Loop")
	for {
		deroHeight := dero.Height()
		moneroHeight := monero.Height()

		if currentDeroHeight != deroHeight {
			currentDeroHeight = deroHeight

			if currentMoneroHeight != moneroHeight {
				currentMoneroHeight = moneroHeight

				moneroTransfers, err := monero.GetIncomingTransfersByHeight(monero.Height())

				// Check if moneroTransfers is empty or has no entries
				if moneroTransfers.Entries.In == nil {
					// Handle the case where there are no transfers

					exports.Logs.Info(dero.Echo("XMR"), "Height:", moneroHeight)
					continue
				}

				if err != nil {
					return err
				}
				if err := processMoneroTransfers(moneroTransfers); err != nil {
					return err
				}
			}

			deroTransfers, err := dero.GetIncomingTransfersByHeight(dero.Height())

			if deroTransfers == nil {
				exports.Logs.Info(dero.Echo("DERO"), "Height:", deroHeight)

				continue
			}

			if err != nil {
				return err
			}

			if err := processDeroTransfers(deroTransfers); err != nil {
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

func Info(message string) {
	msg := dero.Echo(message)
	exports.Logs.Info(msg)
}
