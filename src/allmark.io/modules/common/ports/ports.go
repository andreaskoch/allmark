// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ports

import (
	"fmt"
	"net"
)

// portRangeStart defines the first port that will be assigned
var portRangeStart = 33000

// portRangeEnd  defines the highest allowed port to be used.
const portRangeEnd = 34000

// GetFreePort returns a free port for the given TCPAddr.
func GetFreePort(network string, baseAddress net.TCPAddr) int {

	for portToTest := portRangeStart; portToTest < portRangeEnd; portToTest++ {

		// Don't reuse an already used port
		// for the next call to GetFreePort
		portRangeStart = portToTest + 1

		// assign a port to the base address
		baseAddress.Port = portToTest

		if isPortFree(network, baseAddress) {
			return portToTest
		}

	}

	panic(fmt.Sprintf("No free port found between %v and %v ", portRangeStart, portRangeEnd))
}

// isPortFree checks if the port for the given network and base address is free or not.
func isPortFree(network string, baseAddress net.TCPAddr) (isFree bool) {
	addr, err := net.ResolveTCPAddr(network, baseAddress.String())
	if err != nil {
		return false
	}
	l, err := net.ListenTCP(network, addr)
	if err != nil {
		return false
	}
	defer l.Close()
	return true
}
