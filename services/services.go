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
	ExpectationPeriod uint64 `json:"period"`      // set if the service is expected to report periodically, format is UnixTime (milli?)
}

// Get returns a service struct if the identifier matches any
// keys in the services bucket. Returns nil if there are no matching buckets
func Get(identifier string) *Service {
	var service Service
	err := database.BucketView(database.BucketKeyServices, func(b *bbolt.Bucket) error {
		sb := b.Bucket([]byte(identifier))
		if sb == nil {
			return errors.New("no bucket found for the given id")
		}
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
			if err != nil {
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
		return nil
	}
	return &service
}

// Delete deletes the given service if it exists
func Delete(identifier string) error {
	return database.BucketUpdate(database.BucketKeyServices, func(b *bbolt.Bucket) error {
		serviceKey := []byte(identifier)
		if b.Bucket(serviceKey) == nil {
			return errors.New("no service for the given id")
		}
		return b.DeleteBucket(serviceKey)
	})
}

// Add adds a new service to monitor
func Add(identifier string, service Service) error {
	return database.BucketUpdate(database.BucketKeyServices, func(b *bbolt.Bucket) error {
		serviceKey := []byte(identifier)
		if b.Bucket(serviceKey) != nil {
			return errors.New("a service has already been registered for the given id")
		}
		// Create the service bucket, sb
		sb, err := b.CreateBucket(serviceKey)
		if err != nil {
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
		return nil
	})
}
