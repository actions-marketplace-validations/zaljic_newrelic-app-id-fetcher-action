package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// This struct is used to unmarshal the JSON returned by the New Relic API.
type Application struct {
	Applications []struct {
		ID                   int       `json:"id"`
		Name                 string    `json:"name"`
		Language             string    `json:"language"`
		HealthStatus         string    `json:"health_status"`
		Reporting            bool      `json:"reporting"`
		LastReportedAt       time.Time `json:"last_reported_at"`
		ResponseTime         float64   `json:"response_time"`
		Throughput           float64   `json:"throughput"`
		ErrorRate            float64   `json:"error_rate"`
		ApdexTarget          float64   `json:"apdex_target"`
		ApdexScore           float64   `json:"apdex_score"`
		HostCount            int       `json:"host_count"`
		InstanceCount        int       `json:"instance_count"`
		AppApdexThreshold    float64   `json:"app_apdex_threshold"`
		EndUserApdexThresh   int       `json:"end_user_apdex_threshold"`
		RealUserMonitoring   bool      `json:"enable_real_user_monitoring"`
		ServerSideConfig     bool      `json:"use_server_side_config"`
		ApplicationServers   []int     `json:"application_servers"`
		Servers              []int     `json:"servers"`
		ApplicationHosts     []int     `json:"application_hosts"`
		ApplicationInstances []int     `json:"application_instances"`
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
func getApplications(newrelicApiKey string, newrelicApiEndpoint string, appName string) (Application, error) {
	// Create a new net/http client.
	client := &http.Client{}

	// Specify data to be sent in the HTTP request body.
	dataString := fmt.Sprintf(`filter[name]=%s`, appName)
	data := strings.NewReader(dataString)

	// Send a HTTP GET request using net/http to the New Relic API endpoint
	// specified in the newrelicApiEndpoint input parameter.
	req, err := http.NewRequest("GET", newrelicApiEndpoint, data)
	if err != nil {
		log.Fatal(err)
	}

	// Set the Api-Key header to the value of the newrelicApiKey input parameter.
	req.Header.Set("Api-Key", newrelicApiKey)

	// Set the Content-Type header to application/x-www-form-urlencoded.
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request using the net/http client.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Close the HTTP response body.
	defer resp.Body.Close()

	// Print http status code.
	fmt.Println(resp.StatusCode)

	// Return an error if the HTTP status code is not 200.
	if resp.StatusCode != 200 {
		return Application{}, errors.New("HTTP status code is not 200")
	}

	// Unmarshal the HTTP response body into the Application struct.
	var applications Application
	err = json.NewDecoder(resp.Body).Decode(&applications)
	if err != nil {
		log.Fatal(err)
	}

	// Return the list of applications.
	return applications, nil
}

// This function returns the application ID of the previously filtered
// application. It is assumed that the application list only contains one
// application.
func getApplicationId(applications Application) int {
	// Return the application ID of the first application in the list of
	// applications.
	return applications.Applications[0].ID
}
