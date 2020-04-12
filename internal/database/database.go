package database

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

// bbolt bucket keys
var (
	BucketKeyConfig        = []byte("config")         // bucket key for config bucket key
	BucketKeyEvents        = []byte("events")         // bucket key for events bucket
	BucketKeyServices      = []byte("services")       // bucket key for services bucket
	BucketKeyServiceTokens = []byte("service-tokens") // bucket key for service-tokens bucket
)

// ErrConfigBucket is returned when bbolt is unable to open the config bucket
// TODO: I'm not sure if this is the idiomatic way to use errors
var ErrConfigBucket = errors.New("unable to open config bucket")

var db *bbolt.DB

// check for fatal errors
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Open opens the database for the given file name or creates it if it doesn't exist.
// Returns true if the database is already configured, false if it was just created.
func Open(dbFileName string) (configured bool) {
	var err error
	db, err = bbolt.Open(dbFileName, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	check(err)
	// Create the necessary bbolt buckets if they don't exist
	db.Update(func(tx *bbolt.Tx) error {
		configured = tx.Bucket(BucketKeyConfig) != nil
		_, err = tx.CreateBucketIfNotExists(BucketKeyConfig)
		check(err)
		_, err = tx.CreateBucketIfNotExists(BucketKeyServices)
		check(err)
		_, err = tx.CreateBucketIfNotExists(BucketKeyEvents)
		check(err)
		_, err = tx.CreateBucketIfNotExists(BucketKeyServiceTokens)
		check(err)
		return nil
	})
	return configured
}

// Close is just a wrapper around db.Close() in order to keep all database
// management in one file
func Close() {
	db.Close()
}

// GetDB gets the database structure
func GetDB() *bbolt.DB {
	return db
}

// BytesToUint64 converts a byte array to a uint64 number, an operation that is
// often repeated for IDs. It is assumed that the data will parse successfully
// (i.e. type checking is performed in an earlier stage).
// If the parsing fails, the function therefore panics
func BytesToUint64(byteData []byte) uint64 {
	uint64Data, err := strconv.ParseUint(string(byteData), 10, 64)
	check(err)
	return uint64Data
}

// Uint64ToBytes converts a uint64 formatted number to a byte array
func Uint64ToBytes(uint64Data uint64) []byte {
	return []byte(strconv.FormatUint(uint64Data, 10))
}

// PrintBucket simply prints the K-V pairs in the bucket
func PrintBucket(eb *bbolt.Bucket) error {
	return eb.ForEach(func(k, v []byte) error {
		fmt.Printf("- %v: %v\n", string(k), string(v))
		return nil
	})
}
