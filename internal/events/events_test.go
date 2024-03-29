package events

import (
	"fmt"
	"os"
	"sermoni/internal/config"
	"sermoni/internal/database"
	"sermoni/internal/services"
	"testing"
)

const serviceToken = "test-service"

func (e1 *Event) equals(e2 *Event) bool {
	switch {
	case e1.ID != e2.ID:
		return false
	case e1.Service != e2.Service:
		return false
	case e1.Status != e2.Status:
		return false
	case e1.Title != e2.Title:
		return false
	case e1.Details != e2.Details:
		return false
	default:
		return true
	}
}

var testEvents = []*Event{
	{
		Service:   1,
		Timestamp: 1586558825515,
		Status:    "ok",
		Title:     "Backup completed successfully",
	},
	{
		Service:   1,
		Timestamp: 1586558838488,
		Status:    "info",
		Title:     "SSH login for user vetle",
		Details:   "User vetle logged in from IP 192.168.10.110",
	},
	{
		Service:   1,
		Timestamp: 1586558848488,
		Status:    "ok",
	},
	{
		Service:   1,
		Timestamp: 1586558949488,
		Status:    "error",
		Title:     "Backup failed",
		Details:   "Backup couldn't complete because the disk is full",
	},
}

func TestAddEvent(t *testing.T) {
	var err error
	for _, event := range testEvents {
		if err = Add(event); err != nil {
			fmt.Println(err)
			t.Fatal("error returned when trying to add event")
		}
	}

	// Assumes that bbolt starts sequences on 1
	for i, event := range testEvents {
		event.ID = uint64(i) + 1
		event.Service = 1
	}
}

func TestGetAll(t *testing.T) {
	events := GetAll()
	for i, event := range events {
		if !event.equals(testEvents[i]) {
			t.Fatal("stored event does not match original")
		}
	}
}

func TestMain(m *testing.M) {
	// (Re)create the test database
	testDB := "test.db"
	os.Remove(testDB)
	database.Open(testDB)
	config.InitConfig()
	err := services.Add(&services.Service{
		Token:       serviceToken,
		Name:        "test @ dev-laptop",
		Description: "Service used for testing only",
	})
	if err != nil {
		print("Couldn't add test service")
		os.Exit(1)
	}
	defer database.Close()
	os.Exit(m.Run())
}
