package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Logger(t *testing.T) {
	type Adapter struct {
		MftVersion string `json:"mft-version"`
		Name       string `json:"name"`
	}

	adapter := Adapter{"1.0", "RAID Controller"}
	json, err := json.Marshal(adapter)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(json))
}
