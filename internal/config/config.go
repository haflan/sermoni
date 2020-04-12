package config

import (
	"log"
	"sermoni/internal/database"

	"github.com/gorilla/securecookie"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

var (
	defaultPassPhrase = []byte("admin")
	defaultPageTitle  = []byte("sermoni")

	keyPassHash   = []byte("passhash")
	keyPageTitle  = []byte("pagetitle")
	keySCHashKey  = []byte("schashkey")  // Secure cookie hash key
	keySCBlockKey = []byte("blockkey")   // Secure cookie block key
	keySessionKey = []byte("sessionkey") // Session key
	keyCSRFKey    = []byte("csrfkey")    // CSRF protection auth key
)

// Config is a struct that contains all configuration parameters as []byte data
type Config struct {
	PassHash   []byte
	PageTitle  []byte
	SCHashKey  []byte
	SCBlockKey []byte
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
			SCHashKey:  b.Get(keySCHashKey),
			SCBlockKey: b.Get(keySCBlockKey),
			SessionKey: b.Get(keySessionKey),
			CSRFKey:    b.Get(keyCSRFKey),
		}
		return nil
	})
	return

}

// InitConfig populates the config root bucket with default configurations
// (Web client) passphrase and page title can be reset later
func InitConfig() {
	db := database.GetDB()
	// TODO: Maybe this belongs elsewhere?
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
	passhash, err := bcrypt.GenerateFromPassword(defaultPassPhrase, bcrypt.DefaultCost)
	hashKey := securecookie.GenerateRandomKey(32)
	blockKey := securecookie.GenerateRandomKey(32)
	sessionKey := securecookie.GenerateRandomKey(32)
	CSRFKey := securecookie.GenerateRandomKey(32)
	check(err)
	db.Update(func(tx *bbolt.Tx) error {
		var err error
		b := tx.Bucket(database.BucketKeyConfig)
		err = b.Put(keyPassHash, passhash)
		check(err)
		err = b.Put(keyPageTitle, defaultPageTitle)
		check(err)
		err = b.Put(keySCHashKey, hashKey)
		check(err)
		err = b.Put(keySCBlockKey, blockKey)
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
