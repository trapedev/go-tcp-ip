package util

import (
	"fmt"
	"testing"
)

func TestGetLocalUIpAddress(t *testing.T) {
	res, err := GetLocalUIpAddress("eth0")
	if err != nil {
		t.Errorf("%v", err)
	} else {
		fmt.Printf("%v", res)
	}
}
