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

var currentDeroTransfers *rpc.Get_Transfers_Result
var currentMoneroTransfers monero.TransferResult

func processIncomingTransfers(db *bbolt.DB, LoopActivated *bool) error {

	processMoneroTransfers := func(transfers monero.TransferResult) error {
		if !*LoopActivated {
			*LoopActivated = true
		}

		for _, e := range transfers.In {
			if err := monero.IncomingTransferEntry(e, db); err != nil {
				return err // Exit on error during initial load
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
		deroTransfers, err := dero.GetIncomingTransfers()
		if err != nil {
			return err
		}

		moneroTransfers, err := monero.GetIncomingTransfers()
		if err != nil {
			return err
		}

		if currentDeroTransfers != nil && len(deroTransfers.Entries) > 0 {
			if len(deroTransfers.Entries) != len(currentDeroTransfers.Entries) {
				go func() error {
					exports.Logs.Info(dero.Echo("DERO new entries found"), "Length:", len(deroTransfers.Entries))
					deroTransfers, _ := dero.GetIncomingTransfersByHeight(dero.Height())
					if err := processDeroTransfers(deroTransfers); err != nil {
						return err
					}
					return nil
				}()
			}
		}

		currentDeroTransfers = &deroTransfers

		if currentMoneroTransfers.In != nil && len(moneroTransfers.In) > 0 {
			if len(moneroTransfers.In) != len(currentMoneroTransfers.In) {
				go func() error {
					exports.Logs.Info(dero.Echo("Monero new entries found"), "Length:", len(moneroTransfers.In), "Height:", monero.Height())
					moneroTransfers, _ = monero.GetIncomingTransfers()
					if err := processMoneroTransfers(moneroTransfers); err != nil {
						return err
					}
					return nil
				}()
			}
		}

		currentMoneroTransfers = moneroTransfers

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
