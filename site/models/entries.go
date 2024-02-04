package models

import "github.com/deroproject/derohe/rpc"

type EntriesData struct {
	Title   string
	Dev     string
	Entries []rpc.Entry
}
