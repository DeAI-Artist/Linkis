package miner

import (
	"github.com/davecgh/go-spew/spew"
	"net"
	"testing"
)

// TestGetLocalIP tests the getLocalIP function
func TestGetLocalIP(t *testing.T) {
	ip, err := GetPublicIP()
	spew.Dump("your current public ip: ", ip)
	if err != nil {
		t.Errorf("getLocalIP returned an unexpected error: %v", err)
	}
	if ip == "" {
		t.Error("getLocalIP returned an empty string")
	}

	// Check if the IP address format is valid
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		t.Errorf("getLocalIP returned a non-IP string: %s", ip)
	}
}
