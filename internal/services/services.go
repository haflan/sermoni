package services

import (
	"errors"
	"log"
	"sermoni/internal/database"
	"strconv"

	"go.etcd.io/bbolt"
)

var (
	keyServiceID          = []byte("id")
	keyServiceName        = []byte("name")
	keyServiceDescription = []byte("description")
	keyServicePeriod      = []byte("period")
)

// Service describes a service that is expected to report
type Service struct {
	ID                uint64 `json:"id"`          // service id, an integer that represents the service
	Name              string `json:"name"`        // service name, usually on the format 'service @ server'
	Description       string `json:"description"` // more detailed description of the service
	ExpectationPeriod uint64 `json:"period"`      // set if the service is expected to report periodically, format is UnixTime (milli?)
}

// GetByToken returns the service structure associated with the token string, if there
// are any matching entries in service-tokens bucket. Returns nil if there are no matches
func GetByToken(token string) *Service {
	id := getIDFromToken(token)
	if id == nil {
		log.Printf("no service found for the token '%v'\n", token)
		return nil
	}
	return get(id)
}

// GetByID returns the service structure associated with the given uint64-formatted
// service ID, if that service exists. Otherwise returns nil
func GetByID(id uint64) *Service {
	return get([]byte(strconv.FormatUint(id, 10)))
}

// GetAll returns all services in the database (TODO)
func GetAll() (services []*Service) {
	db := database.GetDB()
	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(database.BucketKeyServices)
		stb := tx.Bucket(database.BucketKeyServiceTokens)
		if b == nil {
			log.Panic("The services bucket does not exist")
		}
		if stb == nil {
			log.Panic("The service-tokens bucket does not exist")
		}

		// Go through all k-v pairs in the service *tokens* bucket, in order to get service bucket IDs
		// Use the ID to get the service bucket and create service fromBucket
		return stb.ForEach(func(_, id []byte) error {
			sb := b.Bucket(id)
			service := new(Service)
			service.fromBucket(id, sb)
			services = append(services, service)
			return nil
		})
	})
	return
}

// Delete deletes the given service if it exists
func Delete(intID uint64) error {
	db := database.GetDB()
	serviceID := []byte(strconv.FormatUint(intID, 10))
	return db.Update(func(tx *bbolt.Tx) error {
		var b, stb *bbolt.Bucket
		if b = tx.Bucket(database.BucketKeyServices); b == nil {
			log.Panic("The services bucket does not exist")
		}
		if stb = tx.Bucket(database.BucketKeyServiceTokens); b == nil {
			log.Panic("The service-tokens bucket does not exist")
		}

		// Delete the entry from root services bucket
		if b.Bucket(serviceID) == nil {
			return errors.New("no service for the given id")
		}
		if err := b.DeleteBucket(serviceID); err != nil {
			return err
		}

		// Find the token entry and delete it from service-tokens bucket
		c := stb.Cursor()
		for token, id := c.First(); token != nil; token, id = c.Next() {
			if string(id) == string(serviceID) {
				return stb.Delete(token)
			}
		}
		return errors.New("service id not found in the service-tokens bucket")

		// TODO: Cascade, i.e. delete all events for the given service
		// Maybe this should be done in the HTTP request handler, though?
	})
}

// Add adds a new service to monitor
// Returns error if the token is unavailable and if the transaction fails in any way
func Add(token string, service *Service) error {
	db := database.GetDB()
	return db.Update(func(tx *bbolt.Tx) error {
		var err error
		var b, sb, stb *bbolt.Bucket
		var serviceIDint uint64
		var serviceID []byte

		// Get the services root bucket
		if b = tx.Bucket(database.BucketKeyServices); b == nil {
			log.Panic("The services bucket does not exist")
		}
		// Get the service-tokens root bucket
		if stb = tx.Bucket(database.BucketKeyServiceTokens); stb == nil {
			log.Panic("The service-tokens bucket does not exist")
		}

		// Check if the service token is available, return error otherwise
		serviceToken := []byte(token)
		if serviceID = stb.Get(serviceToken); serviceID != nil {
			return errors.New("a service has already been registered for the given token")
		}

		// Create a new service bucket, sb, and populate it with data from service
		if serviceIDint, err = b.NextSequence(); err != nil {
			return err
		}
		serviceID = []byte(strconv.FormatUint(serviceIDint, 10))
		if sb, err = b.CreateBucket(serviceID); err != nil {
			return err
		}
		if err = sb.Put(keyServiceName, []byte(service.Name)); err != nil {
			return err
		}
		if err = sb.Put(keyServiceDescription, []byte(service.Description)); err != nil {
			return err
		}
		periodStr := strconv.FormatUint(service.ExpectationPeriod, 10)
		if err = sb.Put(keyServicePeriod, []byte(periodStr)); err != nil {
			return err
		}

		// Put an entry in the service-tokens bucket to map the token to the service
		return stb.Put([]byte(token), serviceID)
	})
}

//
// Package-local helpers
//

// fromBucket populates the service struct with data from the given service bucket
// TODO: Consider failing on missing fields and generally choosing an approach more similar to Event.fromBucket
func (service *Service) fromBucket(id []byte, sb *bbolt.Bucket) {
	// Ignoring this error, because it shouldn't be possible
	idInt := database.BytesToUint64(id)
	service.ID = idInt
	if name := sb.Get(keyServiceName); name != nil {
		service.Name = string(name)
	}
	if description := sb.Get(keyServiceDescription); description != nil {
		service.Description = string(description)
	}
	if period := sb.Get(keyServicePeriod); period != nil {
		// Quick fix: Convert to string, then int
		// Uses default value 0 if an error occurs
		intPeriod, err := strconv.ParseUint(string(period), 10, 64)
		if err == nil {
			service.ExpectationPeriod = intPeriod
		} else {
			log.Println("Couldn't convert period to int for service")
			service.ExpectationPeriod = 0
		}
	}
}

// get returns the service structure associated with the []byte-formatted service ID
func get(id []byte) *Service {
	var service Service
	db := database.GetDB()
	err := db.View(func(tx *bbolt.Tx) error {

		// Get the root services bucket and the requested service bucket
		var b, sb *bbolt.Bucket
		if b = tx.Bucket(database.BucketKeyServices); b == nil {
			log.Panic("The services bucket does not exist")
		}
		if sb = b.Bucket(id); sb == nil {
			return errors.New("no service found for the given id")
		}

		// Get service information from the bucket
		service.fromBucket(id, sb)
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return &service
}

// getIDFromToken looks up the given token in the service-tokens bucket and returns
// the ID if it's found, otherwise returning nil
func getIDFromToken(token string) []byte {
	var id []byte
	db := database.GetDB()
	db.View(func(tx *bbolt.Tx) error {
		stb := tx.Bucket(database.BucketKeyServiceTokens)
		if stb == nil {
			log.Panic("The service-tokens bucket does not exist")
		}
		id = stb.Get([]byte(token))
		return nil
	})
	return id
}
