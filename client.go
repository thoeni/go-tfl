package tfl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Client is implemented as default client or cached client
type Client interface {
	GetTubeStatus() ([]Report, error)
	SetBaseURL(newURL string)
}

type DefaultClient struct {
	http    *http.Client
	baseURL string
}

// SetBaseURL sets a custom URL if the default TFL one needs to be changed
func (c *DefaultClient) SetBaseURL(newURL string) {
	c.baseURL = newURL
}

// GetTubeStatus retrieves Tube status
func (c *DefaultClient) GetTubeStatus() ([]Report, error) {
	url := c.baseURL + "Line/Mode/tube,dlr,overground,tflrail/Status/"

	res, err := c.http.Get(url)
	if err != nil {
		log.Print("Couldn't get TFL data")
		log.Print(err)
		return nil, err
	}
	defer res.Body.Close()

	return decodeTflResponse(res.Body)
}

// NewClient returns a pointer to a TFL client
func NewClient(c *http.Client) *DefaultClient {
	client := DefaultClient{
		http:    c,
		baseURL: "https://api.tfl.gov.uk/",
	}
	return &client
}

// InMemoryCachedClient embeds a Client and caches the result
type InMemoryCachedClient struct {
	Client                      Client
	TubeStatus                  []Report
	LastUpdated                 time.Time
	InvalidateIntervalInSeconds float64
}

// GetTubeStatus retrieves Tube status if cache has expired and saves the result back into the cache
func (c *InMemoryCachedClient) GetTubeStatus() ([]Report, error) {
	if time.Since(c.LastUpdated).Seconds() > c.InvalidateIntervalInSeconds {
		r, e := c.Client.GetTubeStatus()
		c.TubeStatus = r
		c.LastUpdated = time.Now()
		return c.TubeStatus, e
	}
	return c.TubeStatus, nil
}

// SetBaseURL sets a custom URL if the default TFL one needs to be changed
func (c *InMemoryCachedClient) SetBaseURL(newURL string) {
	c.Client.SetBaseURL(newURL)
}

// NewCachedClient returns a pointer to a TFL in memory cached client
func NewCachedClient(c *http.Client, cacheDurationSeconds int) *InMemoryCachedClient {
	client := InMemoryCachedClient{
		Client:                      NewClient(c),
		TubeStatus:                  []Report{},
		LastUpdated:                 time.Now().Add(-time.Duration(cacheDurationSeconds+1) * time.Second),
		InvalidateIntervalInSeconds: float64(cacheDurationSeconds),
	}
	return &client
}

func decodeTflResponse(resp io.Reader) ([]Report, error) {
	decoder := json.NewDecoder(resp)

	var data []Report
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return data, nil
}

// ReportArrayToMap helps to convert the []Report into a map[string]Report for easier accessibility
// The key will be the line name toLower(case)
func ReportArrayToMap(reportArray []Report) map[string]Report {
	reportMap := make(map[string]Report)
	for _, report := range reportArray {
		reportMap[strings.ToLower(report.Name)] = report
	}
	return reportMap
}
