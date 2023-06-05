package biz

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	var (
		bts1 []byte
		err  error
		raw  json.RawMessage
	)

	//
	bts1 = []byte(`{"hello": "world"}`)
	err = json.Unmarshal(bts1, &raw)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("~~~ raw:", raw)

	//
	bts1 = []byte(`{"hello": "world"`)
	err = json.Unmarshal(bts1, &raw)
	if err == nil {
		t.Fatal("err shouldn't be nil")
	}
}
