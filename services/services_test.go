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
		t.Fatal("unexpected error when adding service")
	}
	if err = Add(token2, Service{Name: "tester2"}); err != nil {
		t.Fatal("unexpected error when adding second service")
	}
	if err = Add(token2, Service{Name: "another tester"}); err == nil {
		t.Fatal("no error returned when trying to re-use a service token")
	}
	err = Add(token3, Service{Name: "third @ tester", ExpectationPeriod: 300003})
	if err != nil {
		t.Fatal("unexpected error when adding third service")
	}

	/*
		// Add new
		fmt.Printf("DELETE ERR: %v\n", services.Delete(
			services.IntID(services.GetIDFromToken("token1"))))
		fmt.Printf("ADD ERR: %v\n", services		fmt.Printf("ADD ERR: %v\n", services.Add("token2", services.Service{Name: "This"}))
		fmt.Printf("ADD ERR: %v\n", services.Add("token2", services.Service{Name: "This again"}))
		fmt.Printf("token1: %+v\n", services.GetByToken("token1"))
		fmt.Printf("DELETE ERR: %v\n", services.Delete(
			services.IntID(services.GetIDFromToken("token1"))))
		fmt.Printf("token1: %+v\n", services.GetByToken("token1"))
		fmt.Printf("DELETE ERR: %v\n", services.Delete(
			services.IntID(services.GetIDFromToken("token1"))))
	*/
}

func TestDeleteService(t *testing.T) {
	return
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
