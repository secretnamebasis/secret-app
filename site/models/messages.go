package models

import "github.com/deroproject/derohe/rpc"

type MessagesData struct {
	Title    string
	Dev      string
	Messages []rpc.Entry
}
