package tfl

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	testDataCorrect   string = "test-data/tflResponse"
	testDataMalformed string = "test-data/tflResponseMalformed"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestDecodeTflResponse(t *testing.T) {
	responseExample, _ := os.OpenFile(testDataCorrect, os.O_RDONLY, 0666)

	statuses, err := decodeTflResponse(responseExample)

	if err != nil {
		t.Error("File could not be unmarshalled into a status array")
	}
	if len(statuses) != 11 {
		t.Error("Unmarshalled the incorrect number of statuses.")
	}
}

func TestDecodeTflResponseMalformed(t *testing.T) {
	malformedResponseExample, _ := os.OpenFile(testDataMalformed, os.O_RDONLY, 0666)

	_, err := decodeTflResponse(malformedResponseExample)
	if err == nil {
		t.Error("File should not be unmarshalled into a status array")
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Error("Client wasn't generated.")
	}
}

func TestGetTubeStatus(t *testing.T) {
	mockTflResponse, _ := ioutil.ReadFile(testDataCorrect)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(mockTflResponse))
	}))
	defer ts.Close()
	client := NewClient()
	client.SetBaseURL(ts.URL + "/")

	statuses, err := client.GetTubeStatus()

	if err != nil {
		t.Error("Client failed to retrieve TFL data from mock server")
	}
	if len(statuses) != 11 {
		t.Error("Client retrieved and unmarshalled an incorrect number of statuses")
	}
}

func TestGetTubeStatusFailure(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, nil)
	}))
	defer ts.Close()
	client := NewClient()
	client.SetBaseURL(ts.URL + "/")

	_, err := client.GetTubeStatus()

	if err == nil {
		t.Error("Client should have failed to retrieve TFL data from mock server")
	}
}
