package functions

import (
	"fmt"
	"os"

	"github.com/deroproject/derohe/globals"
	"go.etcd.io/bbolt"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateBucket(db *bbolt.DB, bucketName []byte) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func CreateDB(db_name string) (*bbolt.DB, error) {

	db, err := bbolt.Open(db_name, 0600, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return db, nil
}

func Logger() error {
	// parse arguments and setup logging, print basic information
	globals.Arguments["--debug"] = true
	exename, err := os.Executable()

	globals.InitializeLog(os.Stdout, &lumberjack.Logger{
		Filename:   exename + ".log",
		MaxSize:    100, // megabytes
		MaxBackups: 2,
	})
	logger = globals.Logger

	return err

}
