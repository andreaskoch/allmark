// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/andreaskoch/allmark/common/config"
	"net"
	"testing"
)

func Test_getURL_IPv4WildcardAddress_URLUsesLocalhost(t *testing.T) {
	// arrange
	endpoint := HTTPEndpoint{
		domain:      "localhost",
		isSecure:    false,
		forceHTTPS:  false,
		tcpBindings: []*config.TCPBinding{},
	}

	inputBinding := config.TCPBinding{
		Network: "tcp4",
		IP:      "0.0.0.0",
		Zone:    "",
		Port:    80,
	}
	expectedURL := "http://localhost:80"

	// act
	url := getURL(endpoint, inputBinding)

	// assert
	if url != expectedURL {
		t.Fail()
		t.Logf("The url for the endpoint %q and the tcp binding %q should be %q but was %q", endpoint.String(), inputBinding.String(), expectedURL, url)
	}
}

func Test_getURL_IPv6WildcardAddress_URLUsesLocalhost(t *testing.T) {
	// arrange
	endpoint := HTTPEndpoint{
		domain:      "localhost",
		isSecure:    false,
		forceHTTPS:  false,
		tcpBindings: []*config.TCPBinding{},
	}

	inputBinding := config.TCPBinding{
		Network: "tcp6",
		IP:      "::",
		Zone:    "",
		Port:    80,
	}
	expectedURL := "http://localhost:80"

	// act
	url := getURL(endpoint, inputBinding)

	// assert
	if url != expectedURL {
		t.Fail()
		t.Logf("The url for the endpoint %q and the tcp binding %q should be %q but was %q", endpoint.String(), inputBinding.String(), expectedURL, url)
	}
}

func Test_isWildcardAddress_IPv4WildcardIsGiven_ResultIsTrue(t *testing.T) {
	// arrange
	inputAddress := net.ParseIP("0.0.0.0")

	// act
	result := isWildcardAddress(inputAddress)

	// assert
	if result != true {
		t.Fail()
	}
}

func Test_isWildcardAddress_IPv6WildcardIsGiven_ResultIsTrue(t *testing.T) {
	// arrange
	inputAddress := net.ParseIP("::")

	// act
	result := isWildcardAddress(inputAddress)

	// assert
	if result != true {
		t.Fail()
	}
}

func Test_isWildcardAddress_NonWildcardAddressesAreGiven_ResultIsFalse(t *testing.T) {
	// arrange
	inputAddresses := []string{
		"127.0.0.1",
		"::1",
		"192.168.2.1",
		"8.8.8.8",
	}

	for _, addressString := range inputAddresses {
		inputAddress := net.ParseIP(addressString)

		// act
		result := isWildcardAddress(inputAddress)

		// assert
		if result != false {
			t.Fail()
			t.Logf("%q is not a wildcard address but isWildcardAddress returned true.", inputAddress.String())
		}
	}
}

func Test_isIPv6Address_IPv6AddressIsGiven_ResultIsTrue(t *testing.T) {
	// arrange
	inputAddresses := []string{
		"::1",
		"::",
		"fe80::9afe:94ff:fe49:e7ba%en0",
		"2003:7a:896f:3b00:9afe:94ff:fe49:e7ba",
	}

	for _, addressString := range inputAddresses {
		inputAddress := net.ParseIP(addressString)

		// act
		result := isIPv6Address(inputAddress)

		// assert
		if result != true {
			t.Fail()
			t.Logf("%q is an IPv6 address but isIPv6Address returned false.", inputAddress.String())
		}
	}
}

func Test_isIPv6Address_IPv4AddressIsGiven_ResultIsFalse(t *testing.T) {
	// arrange
	inputAddresses := []string{
		"0.0.0.0",
		"127.0.0.1",
		"192.168.2.2",
		"8.8.8.8",
	}

	for _, addressString := range inputAddresses {
		inputAddress := net.ParseIP(addressString)

		// act
		result := isIPv6Address(inputAddress)

		// assert
		if result != false {
			t.Fail()
			t.Logf("%q is an IPv4 address but isIPv6Address returned true.", inputAddress.String())
		}
	}
}
