package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

func RequestToAPI(url, method string, reqHeader map[string]string, reqBody map[string]interface{}, timeOutSecond int) ([]byte, error) {
	timeout := time.Duration(timeOutSecond) * time.Second

	// Convert reqBody to io.Reader
	payload := io.Reader(nil)
	var err error

	if reqBody != nil {
		payload, err = MapToJSONReader(reqBody)
		if err != nil {
			log.Errorf("RequestToAPI 1: " + err.Error())
			return nil, err
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Errorf("RequestToAPI 2: " + err.Error())
		return nil, err
	}

	// Add custom header
	if reqHeader != nil {
		for key, value := range reqHeader {
			if key == "Host" || key == "host" {
				req.Host = value
			} else {
				req.Header.Set(key, value)
			}
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("RequestToAPI 3: " + err.Error())
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return nil, err
	}
	return body, nil
}

func MapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	// Convert the map to a JSON byte slice
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error marshalling map to JSON: %v", err)
	}

	// Create an io.Reader from the JSON byte slice
	return bytes.NewReader(jsonData), nil
}

func IsNodeInSlice(node string, nodes []string) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func IsDomainFormat(domain string) bool {
	// Regular expression to validate domain format
	// - Starts and ends with alphanumeric characters
	// - Allows hyphens but not consecutively or at the start/end
	// - Contains a valid TLD (e.g., .com, .net, .org)
	// - Supports subdomains
	regex := `^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`

	re := regexp.MustCompile(regex)
	return re.MatchString(domain)
}
