package events

import (
	"errors"
	"fmt"
	"log"
	"sermoni/internal/database"
	"strconv"

	"go.etcd.io/bbolt"
)

var (
	keyEventID        = []byte("id")
	keyEventService   = []byte("service")
	keyEventTimestamp = []byte("timestamp")
	keyEventStatus    = []byte("status")
	keyEventTitle     = []byte("title")
	keyEventDetails   = []byte("details")
)

// Event contains data sent to sermoni from a service
type Event struct {
	ID        uint64 `json:"id"`
	Service   uint64 `json:"service"` // ID of a Service (to be mapped to service name client-side)
	Timestamp uint64 `json:"timestamp"`
	Status    string `json:"status"`
	Title     string `json:"title"`
	Details   string `json:"details"`
}

// GetAll returns all events in the database
func GetAll() (events []*Event) {
	db := database.GetDB()
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(database.BucketKeyEvents)
		if b == nil {
			log.Panic("The events bucket does not exist")
		}
		// ForEach doesn't return buckets (nil instead), so only the key is useful
		return b.ForEach(func(id, _ []byte) error {
			eb := b.Bucket(id)
			event := new(Event)
			if err := event.fromBucket(eb); err != nil {
				return err
			}
			events = append(events, event)
			return nil
		})
	})
	if err != nil {
		fmt.Println(err)
	}
	return
}

// Delete event with the given ID.
// Returns error if no such event can be found
func Delete(idInt uint64) error {
	db := database.GetDB()
	id := database.Uint64ToBytes(idInt)
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(database.BucketKeyEvents)
		if b == nil {
			log.Panic("The events bucket does notexist")
		}
		return b.DeleteBucket(id)
	})
}

// Add persists a new event to database after generating an ID for it
func Add(event *Event) error {
	db := database.GetDB()
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(database.BucketKeyEvents)
		if b == nil {
			log.Panic("The events bucket does not exist")
		}

		// Create a new event ID
		idInt, err := b.NextSequence()
		if err != nil {
			return err
		}
		id := []byte(strconv.FormatUint(idInt, 10))
		event.ID = idInt

		// Create the event bucket and fill it with data from event
		eb, err := b.CreateBucket(id)
		if err != nil {
			return err
		}
		return event.toBucket(eb)
	})
}

// Writes the event data to the given bucket
func (event *Event) toBucket(eb *bbolt.Bucket) error {
	var err error
	id := database.Uint64ToBytes(event.ID)
	serviceID := database.Uint64ToBytes(event.Service)
	timestamp := database.Uint64ToBytes(event.Timestamp)
	if err = eb.Put(keyEventID, id); err != nil {
		return err
	}
	if eb.Put(keyEventService, serviceID); err != nil {
		return err
	}
	if eb.Put(keyEventTimestamp, timestamp); err != nil {
		return err
	}
	if eb.Put(keyEventStatus, []byte(event.Status)); err != nil {
		return err
	}
	if eb.Put(keyEventTitle, []byte(event.Title)); err != nil {
		return err
	}
	if eb.Put(keyEventDetails, []byte(event.Details)); err != nil {
		return err
	}
	return nil
}

// Reads data from the given bucket into the fields of event
// Returns error if any of the fields cannot be found
func (event *Event) fromBucket(eb *bbolt.Bucket) error {
	var id, service, timestamp []byte
	var status, title, details []byte
	err := errors.New("missing field from database")

	// Get data from database
	if id = eb.Get(keyEventID); id == nil {
		return err
	}
	if service = eb.Get(keyEventService); service == nil {
		return err
	}
	if timestamp = eb.Get(keyEventTimestamp); timestamp == nil {
		return err
	}
	if status = eb.Get(keyEventStatus); status == nil {
		return err
	}
	if title = eb.Get(keyEventTitle); title == nil {
		return err
	}
	if details = eb.Get(keyEventDetails); details == nil {
		return err
	}

	// Format data and set fields of event
	event.ID = database.BytesToUint64(id)
	event.Service = database.BytesToUint64(service)
	event.Timestamp = database.BytesToUint64(timestamp)
	event.Status = string(status)
	event.Title = string(title)
	event.Details = string(details)

	return nil
}
