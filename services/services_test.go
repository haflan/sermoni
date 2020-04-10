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

func TestAddService(t *testing.T) {
	token1 := "my-great-token"
	token2 := "my-other-token"
	token3 := "my-third-token"
	err := Add(token1, Service{
		Name:              "tester @ dev-computer",
		Description:       "This describes the service in more detail",
		ExpectationPeriod: 282342,
	})
	if err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when adding service")
	}
	if err = Add(token2, Service{Name: "tester2"}); err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when adding second service")
	}
	if err = Add(token2, Service{Name: "another tester"}); err == nil {
		t.Fatal("no error returned when trying to re-use a service token")
	}
	err = Add(token3, Service{Name: "third @ tester", ExpectationPeriod: 300003})
	if err != nil {
		t.Fatal("unexpected error when adding third service")
	}
}


func TestDeleteService(t *testing.T) {
	// bbolt should always start ID sequences on 1, so hard coding IDs should be fine
	err := Delete(2)
	if err != nil {
		fmt.Println(err)
		t.Fatal("unexpected error when trying to delete service")
	}
	if err = Delete(2); err == nil {
		t.Fatal("no error returned when trying to delete non-existing service")
	}
}

func TestGetAll(t *testing.T) {
	services := GetAll()
	for _, service := range services {
		fmt.Printf("%+v\n", service)
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
