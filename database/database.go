package database

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"go.etcd.io/bbolt"
)

const (
	defaultPassPhrase = "admin"
	defaultPageTitle  = "sermoni"
)

// bbolt keys
var (
	BucketKeyConfig   = []byte("config")   // bucket key for config bucket key
	BucketKeyServices = []byte("services") // bucket key for services bucket
	BucketKeyEvents   = []byte("events")   // bucket key for events bucket

	keyPassHash  = []byte("passhash")
	keyPageTitle = []byte("pagetitle")
)

// ErrConfigBucket is returned when bbolt is unable to open the config bucket
// TODO: I'm not sure if this is the idiomatic way to use errors
var ErrConfigBucket = errors.New("unable to open config bucket")

var db *bbolt.DB

// Init opens the database for the given file name or creates it if it doesn't exist.
// It also populates it with essential configuration data if required.
func Init(dbFileName string) error {
	fmt.Printf("Init db '%v'\n", dbFileName)
	var err error
	db, err = bbolt.Open(dbFileName, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	// Create the necessary bbolt buckets if they don't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(BucketKeyConfig); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(BucketKeyServices); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(BucketKeyEvents); err != nil {
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
		b := tx.Bucket(BucketKeyConfig)
		if b == nil {
			return ErrConfigBucket
		}
		passhash := b.Get(keyPassHash)
		configured = passhash != nil
		return nil
	})
	if !configured {
		return Reconfigure(defaultPassPhrase, defaultPageTitle)
	}
	return nil
}

// Reconfigure takes a passphrase and a page title for the web page,
// generates hash for the password and updates the database with this
// new configuration.
func Reconfigure(passphrase string, pageTitle string) error {

	/* TODO: Generate a random _readable_ password if none is given
	var passphraseBytes []byte
	if passphrase == "" {
		passphraseBytes = make([]byte, 24)
		rand.Read(passphraseBytes)
		passphrase = string(passphraseBytes)
		fmt.Printf("Generated passphrase: %v\n", []rune(passphrase))
	} else {
		passphraseBytes = []byte(passphrase)
	}*/
	//sha256.Sum256([]byte(passphraseBytes))

	// TODO: Maybe bcrypt is overkill for such a small project? Consider later
	passhash, err := bcrypt.GenerateFromPassword([]byte(passphrase), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(BucketKeyConfig)
		if b == nil {
			return ErrConfigBucket
		}
		var err error
		if err = b.Put(keyPassHash, passhash); err != nil {
			return err
		}
		if err = b.Put(keyPageTitle, []byte(pageTitle)); err != nil {
			return err
		}
		return nil
	})
}

// Close is just a wrapper around db.Close() in order to keep all database
// management in one file
func Close() {
	db.Close()
}

// GetDB returns the bbolt database struct
// TODO: Not sure how if this is a good way to do it, although it does seem to work
func GetDB() *bbolt.DB {
	return db
}
