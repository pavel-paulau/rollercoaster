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

func generateId() io.Reader {
	i, _ := json.Marshal(struct {
		ID uint64 `json:"id"`
	}{
		ID: 1,
	})

	return bytes.NewReader(i)
}

func TestOpenDB(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("open did not panic")
		}
	}()

	dbName = "bad/name"
	db = open()
}

func TestPost(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

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
		t.Fatalf("expected: 201, got: %d", resp.StatusCode)
	}

	if msg.Message != "ok" {
		t.Errorf("expected: ok, got: %s", msg.Message)
	}
}

func TestBadPayload(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

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
		t.Fatalf("expected: 400, got: %d", resp.StatusCode)
	}

	if msg.Message != "EOF" {
		t.Errorf("expected: EOF, got: %s", msg.Message)
	}
}

func TestPostDBError(t *testing.T) {
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

	if resp.StatusCode != 500 {
		t.Fatalf("expected: 500, got: %d", resp.StatusCode)
	}

	if msg.Message != "database not open" {
		t.Errorf("expected: 'database not open', got: %s", msg.Message)
	}
}

func TestGet(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

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
		t.Fatalf("expected: 200, got: %d", resp.StatusCode)
	}

	if len(benchmarks) != 1 {
		t.Fatalf("expected: 1 benchmark, got: %d", len(benchmarks))
	}

	if benchmarks[0].Group != "myGroup" {
		t.Fatalf("expected: %s, got: %s", "myGroup", benchmarks[1].Group)
	}
}

func TestDelete(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	http.Post(ts.URL+"/api/v1/benchmarks", "application/json", generateBenchmark())

	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/benchmarks", generateId())
	resp, err := http.DefaultClient.Do(req)
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

	if resp.StatusCode != 200 {
		t.Fatalf("expected: 201, got: %d", resp.StatusCode)
	}

	if msg.Message != "ok" {
		t.Errorf("expected: ok, got: %s", msg.Message)
	}
}

func TestEmptyDelete(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()
	defer db.Close()

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	http.Post(ts.URL+"/api/v1/benchmarks", "application/json", generateBenchmark())

	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/benchmarks", nil)
	resp, err := http.DefaultClient.Do(req)
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
		t.Fatalf("expected: 201, got: %d", resp.StatusCode)
	}

	if msg.Message != "EOF" {
		t.Errorf("expected: EOF, got: %s", msg.Message)
	}
}

func TestDeleteDBError(t *testing.T) {
	tmp, _ := ioutil.TempFile("", "rollercoaster")
	defer os.Remove(tmp.Name())

	dbName = tmp.Name()
	db = open()
	initBucket()

	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	http.Post(ts.URL+"/api/v1/benchmarks", "application/json", generateBenchmark())
	db.Close()

	req, err := http.NewRequest("DELETE", ts.URL+"/api/v1/benchmarks", generateId())
	resp, err := http.DefaultClient.Do(req)
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

	if resp.StatusCode != 500 {
		t.Fatalf("expected: 500, got: %d", resp.StatusCode)
	}

	if msg.Message != "database not open" {
		t.Errorf("expected: 'database not open', got: %s", msg.Message)
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
		t.Fatalf("expected: 200, got: %d", resp.StatusCode)
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
		t.Fatalf("expected: 200, got: %d", resp.StatusCode)
	}
}
