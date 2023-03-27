package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// Here we define the struct that will be used to unmarshal the JSON response
// from the New Relic API.
type NewRelicApplications struct {
	Applications []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Language     string `json:"language"`
		HealthStatus string `json:"health_status"`
		Reporting    bool   `json:"reporting"`
		Settings     struct {
			AppApdexThreshold        float64 `json:"app_apdex_threshold"`
			EndUserApdexThreshold    int     `json:"end_user_apdex_threshold"`
			EnableRealUserMonitoring bool    `json:"enable_real_user_monitoring"`
			UseServerSideConfig      bool    `json:"use_server_side_config"`
		} `json:"settings"`
		Links struct {
			ApplicationInstances []any `json:"application_instances"`
			Servers              []any `json:"servers"`
			ApplicationHosts     []any `json:"application_hosts"`
		} `json:"links"`
		LastReportedAt     time.Time `json:"last_reported_at,omitempty"`
		ApplicationSummary struct {
			ResponseTime  float64 `json:"response_time"`
			Throughput    int     `json:"throughput"`
			ErrorRate     float64 `json:"error_rate"`
			ApdexTarget   float64 `json:"apdex_target"`
			ApdexScore    float64 `json:"apdex_score"`
			HostCount     int     `json:"host_count"`
			InstanceCount int     `json:"instance_count"`
		} `json:"application_summary,omitempty"`
	} `json:"applications"`
	Links struct {
		ApplicationServers              string `json:"application.servers"`
		ApplicationServer               string `json:"application.server"`
		ApplicationApplicationHosts     string `json:"application.application_hosts"`
		ApplicationApplicationHost      string `json:"application.application_host"`
		ApplicationApplicationInstances string `json:"application.application_instances"`
		ApplicationApplicationInstance  string `json:"application.application_instance"`
	} `json:"links"`
}

// This function is the entry point for the action. It is responsible for
// parsing the input parameters, calling the functions that fetch the
// application ID from the New Relic API, and setting the output parameter.
func main() {
	// Get the input parameters from the environment variables.
	newrelicApiKey := os.Getenv("INPUT_NEWRELICAPIKEY")
	newrelicRegion := os.Getenv("INPUT_NEWRELICREGION")
	appName := os.Getenv("INPUT_APPNAME")

	// Return an error if the newrelicApiKey input parameter is not set.
	if newrelicApiKey == "" {
		fmt.Println("NewRelic API key not specified.")
		os.Exit(1)
	}

	// Return an error if the appName input parameter is not set.
	if appName == "" {
		fmt.Println("App name not specified.")
		os.Exit(1)
	}

	// Set the New Relic API endpoint based on the region specified in the
	// newrelicRegion input parameter.
	newrelicApiEndpoint := ""
	// The New Relic API endpoint is different for US and EU regions.
	if newrelicRegion == "US" {
		newrelicApiEndpoint = "https://api.newrelic.com/v2/applications.json"
	} else if newrelicRegion == "EU" {
		newrelicApiEndpoint = "https://api.eu.newrelic.com/v2/applications.json"
		// If the region is not US or EU, exit with an error.
	} else {
		fmt.Println("Invalid NewRelic region specified.")
		os.Exit(1)
	}

	// Call the getApplications function to fetch the list of applications from
	// the New Relic API.
	applications, err := getApplications(newrelicApiKey, newrelicApiEndpoint, appName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Call the getApplicationId function to get the application ID from the
	// filtered list of applications.
	applicationId := getApplicationId(applications)

	// Set the output parameter to the application ID.
	output := fmt.Sprintf("%d", applicationId)
	// Print the output parameter to stdout.
	fmt.Printf(`::set-output name=appID::%s`, output)
}

// This function sends a HTTP GET request to the endpoint specified in the
// newrelicApiEndpoint input parameter and returns the list of applications
// returned by the New Relic API. It filters the list of applications to only
// include applications with the name specified in the appName input
// parameter.
func getApplications(newrelicApiKey string, newrelicApiEndpoint string, appName string) ([]NewRelicApplications, error) {
	// Specify data to be sent in the HTTP request body.
	dataString := fmt.Sprintf(`filter[name]=%s`, appName)
	data := strings.NewReader(dataString)

	// Send a HTTP GET request using net/http to the New Relic API endpoint
	// specified in the newrelicApiEndpoint input parameter.
	req, err := http.NewRequest("GET", newrelicApiEndpoint, data)
	if err != nil {
		return nil, err
	}

	// Set the Api-Key header to the value of the newrelicApiKey input parameter.
	req.Header.Set("Api-Key", newrelicApiKey)

	// Send the HTTP request using the net/http client.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Close the HTTP response body.
	defer resp.Body.Close()

	// Unmarshal the JSON response into a Application struct.
	var applications []NewRelicApplications
	err = json.NewDecoder(resp.Body).Decode(&applications)
	if err != nil {
		return nil, err
	}

	// Return the list of applications. If the list is empty, return an error.
	if len(applications) > 0 {
		return applications, nil
	} else {
		return nil, fmt.Errorf("no applications found with the name %s", appName)
	}

}

// This function returns the application ID of the previously filtered
// application. It is assumed that the application list only contains one
// application.
func getApplicationId(applications []NewRelicApplications) int {
	// Return the application ID of the first application in the list of
	// applications.
	return applications[0].Applications[0].ID
}
