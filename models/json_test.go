package models

import (
	_ "embed"
	// "fmt"
	"encoding/json"
	"testing"
)

var (
	//go:embed data.json
	_DataBytes []byte
)

func BenchmarkUnmarshal_01(b *testing.B) {
	b.ReportAllocs()

	type Data struct {
		Data any `json:"data"`
	}

	for i := 0; i < b.N; i++ {
		d := new(Data)
		e := json.Unmarshal(_DataBytes, d)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkMarshal_01(b *testing.B) {
	b.ReportAllocs()

	type Data struct {
		Data any `json:"data"`
	}

	b.StopTimer()
	d := new(Data)
	e := json.Unmarshal(_DataBytes, d)
	if e != nil {
		b.Fatal(e)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, e := json.Marshal(d)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkUnmarshal_02(b *testing.B) {
	b.ReportAllocs()

	type Data struct {
		Data json.RawMessage `json:"data"`
	}

	for i := 0; i < b.N; i++ {
		d := new(Data)
		e := json.Unmarshal(_DataBytes, d)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkMarshal_02(b *testing.B) {
	b.ReportAllocs()

	type Data struct {
		Data json.RawMessage `json:"data"`
	}

	b.StopTimer()
	d := new(Data)
	e := json.Unmarshal(_DataBytes, d)
	if e != nil {
		b.Fatal(e)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, e := json.Marshal(d)
		if e != nil {
			b.Fatal(e)
		}
	}
}
