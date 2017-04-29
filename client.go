package tfl

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Client holds few properties and is receiver for few methods to interact with TFL apis
type Client struct {
	baseURL string
}

// NewClient returns a pointer to a TFL client
func NewClient() *Client {
	client := Client{
		baseURL: "https://api.tfl.gov.uk/",
	}
	return &client
}

// SetBaseURL sets a custom URL if the default TFL one needs to be changed
func (c *Client) SetBaseURL(newURL string) {
	c.baseURL = newURL
}

// GetTubeStatus retrieves Tube status
func (c *Client) GetTubeStatus() ([]Report, error) {
	url := c.baseURL + "Line/Mode/tube,dlr,overground,tflrail/Status/"

	res, err := http.Get(url)
	if err != nil {
		log.Print("Couldn't get TFL data")
		log.Print(err)
		return nil, err
	}
	defer res.Body.Close()

	return decodeTflResponse(res.Body)
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