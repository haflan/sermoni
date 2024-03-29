package config

import (
	"crypto/sha256"
	"log"
	"sermoni/internal/database"

	"github.com/gorilla/securecookie"
	"go.etcd.io/bbolt"
)

var (
	defaultPassPhrase = []byte("admin")
	defaultPageTitle  = []byte("sermoni")

	keyPassHash   = []byte("passhash")
	keyPageTitle  = []byte("pagetitle")
	keySessionKey = []byte("sessionkey") // Session key
	keyCSRFKey    = []byte("csrfkey")    // CSRF protection auth key
)

// Config is a struct that contains all configuration parameters as []byte data
type Config struct {
	PassHash   []byte
	PageTitle  []byte
	SessionKey []byte
	CSRFKey    []byte
}

// GetConfig Creates a Config struct from the values in database
// Should only be necessary to call once
func GetConfig() (config *Config) {
	db := database.GetDB()
	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(database.BucketKeyConfig)
		config = &Config{
			PassHash:   b.Get(keyPassHash),
			PageTitle:  b.Get(keyPageTitle),
			SessionKey: b.Get(keySessionKey),
			CSRFKey:    b.Get(keyCSRFKey),
		}
		return nil
	})
	return

}

// SetPassphrase persists the sha256sum of the given passphrase to DB
func SetPassphrase(passphrase string) {
	db := database.GetDB()
	passhash := sha256.Sum256([]byte(passphrase))
	db.Update(func(tx *bbolt.Tx) error {
		var err error
		b := tx.Bucket(database.BucketKeyConfig)
		err = b.Put(keyPassHash, passhash[:])
		check(err)
		return nil
	})
}

// SetPageTitle persists the given page title to DB
func SetPageTitle(pageTitle string) {
	db := database.GetDB()
	db.Update(func(tx *bbolt.Tx) error {
		var err error
		b := tx.Bucket(database.BucketKeyConfig)
		err = b.Put(keyPageTitle, []byte(pageTitle))
		check(err)
		return nil
	})
}

// InitConfig populates the config root bucket with default configurations
// (Web client) passphrase and page title can be reset later
func InitConfig() {
	db := database.GetDB()
	// TODO: Generate a random _readable_ password if none is given
	passhash := sha256.Sum256([]byte(defaultPassPhrase))
	sessionKey := securecookie.GenerateRandomKey(32)
	CSRFKey := securecookie.GenerateRandomKey(32)
	db.Update(func(tx *bbolt.Tx) error {
		var err error
		b := tx.Bucket(database.BucketKeyConfig)
		err = b.Put(keyPassHash, passhash[:])
		check(err)
		err = b.Put(keyPageTitle, defaultPageTitle)
		check(err)
		err = b.Put(keySessionKey, sessionKey)
		check(err)
		err = b.Put(keyCSRFKey, CSRFKey)
		check(err)
		return nil
	})
}

// check for fatal errors
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
