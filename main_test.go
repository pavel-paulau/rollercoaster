package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func generateBenchmark() io.Reader {
	b, _ := json.Marshal(benchmark{
		Group:  "myGroup",
		Metric: "myMetric, ops/sec",
		Value:  1.23,
	})

	return bytes.NewReader(b)
}

func TestPost(t *testing.T) {
	tmp, err := ioutil.TempFile("", "rollercoaster")

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()
	defer os.Remove(tmp.Name())

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/v1/benchmarks", "application/json", generateBenchmark())
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	msg := struct {
		Message string `json:"message"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("Expected: 201, got: %d", resp.StatusCode)
	}

	if msg.Message != "ok" {
		t.Errorf("Expected: ok, got: %s", msg.Message)
	}
}

func TestBadPayload(t *testing.T) {
	tmp, err := ioutil.TempFile("", "rollercoaster")

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()
	defer os.Remove(tmp.Name())

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/v1/benchmarks", "application/json", nil)
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	msg := struct {
		Message string `json:"message"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 400 {
		t.Fatalf("Expected: 400, got: %d", resp.StatusCode)
	}

	if msg.Message != "EOF" {
		t.Errorf("Expected: EOF, got: %s", msg.Message)
	}
}

func TestGet(t *testing.T) {
	tmp, err := ioutil.TempFile("", "rollercoaster")

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()
	defer os.Remove(tmp.Name())

	http.Post(ts.URL+"/api/v1/benchmarks", "application/json", generateBenchmark())

	resp, err := http.Get(ts.URL + "/api/v1/benchmarks")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	var benchmarks []benchmark
	err = json.NewDecoder(resp.Body).Decode(&benchmarks)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected: 200, got: %d", resp.StatusCode)
	}

	if len(benchmarks) != 1 {
		t.Fatalf("Expected: 1 benchmark, got: %d", len(benchmarks))
	}

	if benchmarks[0].Group != "myGroup" {
		t.Fatalf("Expected: %s, got: %s", "myGroup", benchmarks[1].Group)
	}
}

func TestMainPage(t *testing.T) {
	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected: 200, got: %d", resp.StatusCode)
	}
}

func TestStaticAssets(t *testing.T) {
	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/static/main.js")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected: 200, got: %d", resp.StatusCode)
	}
}
