package services

import (
	"errors"
	"log"
	"sermoni/database"
	"strconv"

	"go.etcd.io/bbolt"
)

var (
	keyServiceName        = []byte("name")
	keyServiceDescription = []byte("description")
	keyServicePeriod      = []byte("period")
)

// Service describes a service that is expected to report
type Service struct {
	Name              string `json:"name"`        // service name, usually on the format 'service @ server'
	Description       string `json:"description"` // more detailed description of the service
	ExpectationPeriod int    `json:"period"`      // set if the service is expected to report periodically, format is UnixTime (milli?)
}

// Get returns a service struct if the identifier matches any
// keys in the services bucket. Returns nil if there are no matching buckets
func Get(identifier string) (service Service) {
	db := database.GetDB()
	err := db.View(func(tx *bbolt.Tx) error {
		serviceBucket := tx.Bucket(database.BucketKeyServices)
		if serviceBucket == nil {
			log.Fatal("No services bucket found")
		}
		b := serviceBucket.Bucket([]byte(identifier))
		if b == nil {
			return errors.New("no bucket found for the given id")
		}
		if name := b.Get(keyServiceName); name != nil {
			service.Name = string(name)
		}
		if description := b.Get(keyServiceDescription); description != nil {
			service.Description = string(description)
		}
		if period := b.Get(keyServicePeriod); period != nil {
			// Quick fix: Convert to string, then int
			// If an error occurs (it tho)
			if intPeriod, err := strconv.Atoi(string(period)); err != nil {
				service.ExpectationPeriod = intPeriod
				log.Printf("Couldn't convert period to int for service with id '%v'\n", identifier)
			} else {
				service.ExpectationPeriod = 0
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return service
}

// Delete deletes the given service if it exists
func Delete(identifier string) {
	//db := database.GetDB()
}

// Add adds a new service to monitor
func Add(identifier string, service Service) error {
	db := database.GetDB()
	return db.Update(func(tx *bbolt.Tx) error {
		serviceBucket := tx.Bucket(database.BucketKeyServices)
		if serviceBucket == nil {
			log.Fatal("No services bucket found")
		}
		serviceKey := []byte(identifier)
		if serviceBucket.Bucket(serviceKey) != nil {
			return errors.New("a service has already been registered for the given id")
		}
		b, err := serviceBucket.CreateBucket(serviceKey)
		if err != nil {
			return err
		}
		if err = b.Put(keyServiceName, []byte(service.Name)); err != nil {
			return err
		}
		if err = b.Put(keyServiceDescription, []byte(service.Description)); err != nil {
			return err
		}
		periodStr := strconv.Itoa(service.ExpectationPeriod)
		if err = b.Put(keyServicePeriod, []byte(periodStr)); err != nil {
			return err
		}
		return nil
	})
}

// TODO: Consider a wrapper that gets the services bucket and operates on it
