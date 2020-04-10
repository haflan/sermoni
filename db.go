package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

const defaultPageTitle = "sermoni"

// bbolt keys
var (
	bucketKeyConfig       = []byte("config")
	bucketKeyServices     = []byte("services")
	bucketKeyEvents       = []byte("events")
	keyPassHash           = []byte("passhash")
	keyPageTitle          = []byte("pagetitle")
	keyServiceName        = []byte("name")
	keyServiceDescription = []byte("description")
	keyServicePeriod      = []byte("period")
)

// ErrConfigBucket is returned when bbolt is unable to open the config bucket
// TODO: I'm not sure if this is the idiomatic way to use errors
var ErrConfigBucket = errors.New("unable to open config bucket")

// Global bbolt DB struct
var db *bbolt.DB

// openDB opens the database for the given file name or creates it if it doesn't exist
func initDB(dbFileName string) error {
	fmt.Printf("Init db '%v'\n", dbFileName)
	var err error
	db, err = bbolt.Open(dbFileName, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// Create the necessary bbolt buckets if they don't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketKeyConfig); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketKeyServices); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketKeyEvents); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Check if the config is initialized - configure if not
	var configured bool
	err = db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKeyConfig)
		if b == nil {
			return ErrConfigBucket
		}
		passhash := b.Get(keyPassHash)
		configured = passhash != nil
		return nil
	})
	if !configured {
		return reconfigure("", defaultPageTitle)
	}
	return nil
}

// reconfigure takes a passphrase and a page title for the web page
// and updates the database with this new configuration.
// If passphrase is the empty string, a random phrase will be generated.
func reconfigure(passphrase string, pageTitle string) error {
	var passphraseBytes []byte
	if passphrase == "" {
		passphraseBytes = make([]byte, 24)
		rand.Read(passphraseBytes)
		passphrase = string(passphraseBytes)
		fmt.Printf("Generated passphrase: %v\n", passphrase)
	} else {
		passphraseBytes = []byte(passphrase)
	}
	passhash := sha256.Sum256(passphraseBytes)
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketKeyConfig)
		if b == nil {
			return ErrConfigBucket
		}
		var err error
		// [:] is needed to get a slice from the [32]byte array
		if err = b.Put(keyPassHash, passhash[:]); err != nil {
			return err
		}
		if err = b.Put(keyPageTitle, []byte(pageTitle)); err != nil {
			return err
		}
		return nil
	})
}

// closeDB is just a wrapper around db.Close() in order to keep all database
// management in one file
func closeDB() {
	db.Close()
}
