package forms

import (
	"os"
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	testConfigInit()
	v, e := GetConfig("list_flags")
	if e != nil {
		fmt.Printf("ERROR - %v\n", e)
	}
	fmt.Printf("Value: %v\n", v)
	if e := SetConfig("NEW_KEY", "New key word 1"); e != nil {
		fmt.Printf("ERROR %v\n", e)
	}
	v, e = GetConfig("config_created")
	fmt.Printf("new value: %v\n",v)
	if e := DeleteConfig("NEW_KEY"); e != nil {
		fmt.Printf("ERROR %v\n", e)
	}
}

func testConfigInit() {
	os.Setenv("DBPATH", "test.db")
	SetupConfigDB()
	if _, e := GetConfig("config_created"); e != nil {
		fmt.Println("Setup default config ....")
		SetupDefaultConfig()
	}
}