package services

import (
	"fmt"
	"os"
	"sermoni/database"
	"strconv"
	"testing"
)

// intID gets uint64 from bytes
func intID(id []byte) uint64 {
	idInt, err := strconv.ParseUint(string(id), 10, 64)
	if err != nil {
		return 0
	}
	return idInt
}

func (s1 *Service) equals (s2 *Service) bool {
	if s1.ID != s2.ID {
		return false
	}
	if s1.Name != s2.Name {
		return false
	}
	if s1.Description != s2.Description {
		return false
	}
	if s1.ExpectationPeriod != s2.ExpectationPeriod {
		return false
	}
	return true
}

var (
	token1 = "my-great-token"
	token2 = "my-other-token"
	token3 = "my-third-token"
)

var testServices = []*Service{
	&Service{
		Name:              "tester @ dev-computer",
		Description:       "This describes the service in more detail",
		ExpectationPeriod: 282342,
	},
	&Service{Name: "tester2", ExpectationPeriod: 300003},
	&Service{Name: "third @ tester"},
}

func TestAddService(t *testing.T) {
	err := Add(token1, testServices[0])
	if err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when adding service")
	}
	if err = Add(token2, testServices[1]); err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when adding second service")
	}
	if err = Add(token2, testServices[1]); err == nil {
		t.Fatal("no error returned when trying to re-use a service token")
	}
	err = Add(token3, testServices[2])
	if err != nil {
		t.Fatal("unexpected error when adding third service")
	}

	// Simulate ID generation for testServices after adding them to DB, to avoid
	// possible interferrence (shouldn't be a problem, but doesn't hurt to be sure).
	// bbolt should always start ID sequences on 1, so this assumes that the service ID 
	// equals the testService index + 1
	for i, service := range testServices {
		service.ID = uint64(i) + 1
	}
}


func TestDeleteService(t *testing.T) {
	var di uint64 = 1	// Deletion index
	err := Delete(di + 1)
	if err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when trying to delete service")
	}
	if err = Delete(di + 1); err == nil {
		t.Fatal("no error returned when trying to delete non-existing service")
	}

	// Delete from testServices too
	testServices = append(testServices[:di], testServices[di+1:]...)
}

func TestGetByID(t *testing.T) {
	var gi uint64 = 0 // Get index
	testService := testServices[gi]
	service := GetByID(gi + 1)
	if !service.equals(testService) {
		t.Fatal("stored service doesn't match original")
	}
}

func TestGetByToken(t *testing.T) {
	testService := testServices[1]
	service := GetByToken(token3)
	if !service.equals(testService) {
		t.Fatal("stored service doesn't match original")
	}
}

func TestGetAll(t *testing.T) {
	services := GetAll()
	for i, service := range services {
		if !service.equals(testServices[i]) {
			t.Fatal("stored service doesn't match original")
		}
	}
}

func TestMain(m *testing.M) {
	// (Re)create the test database
	testDB := "test.db"
	os.Remove(testDB)
	if err := database.Init(testDB); err != nil {
		print("Couldn't initialize test database")
		os.Exit(1)
	}
	defer database.Close()
	os.Exit(m.Run())
}
