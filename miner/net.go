package miner

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetPublicIP() (string, error) {
	response, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", fmt.Errorf("failed to get public IP: %v", err)
	}
	fmt.Println("HTTP Status Code:", response.StatusCode) // Print status code for debugging
	defer response.Body.Close()
	ip, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	return string(ip), nil
}
