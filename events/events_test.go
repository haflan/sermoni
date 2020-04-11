package events

import (
	"os"
	"sermoni/database"
	"sermoni/services"
	"testing"
)

const serviceToken = "test-service"

func TestAddEvent(t *testing.T) {
	Add(serviceToken, &Event{
		Timestamp: 1586558825515,
		Title:     "Backup completed successfully",
	})
}

func TestMain(m *testing.M) {
	// (Re)create the test database
	testDB := "test.db"
	os.Remove(testDB)
	var err error
	if err = database.Init(testDB); err != nil {
		print("Couldn't initialize test database")
		os.Exit(1)
	}
	err = services.Add(serviceToken, &services.Service{
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
