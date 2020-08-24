// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package util

import (
	"fmt"
	"net"
	"runtime"

	"github.com/joeshaw/multierror"
)

// GetNetInfo returns lists of IPs and MACs for the machine it is executed on.
func GetNetInfo() (ipList []string, hwList []string,maskList []string, err error) {
	// Get all interfaces and loop through them
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, nil,nil, err
	}
	sysType := runtime.GOOS
	// Keep track of all errors
	var errs multierror.Errors
    Getaway()
	for _, i := range ifaces {
		// Skip loopback interfaces
		if i.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}

		hw := i.HardwareAddr.String()
		// Skip empty hardware addresses
		if hw != "" {
			hwList = append(hwList, hw)
		}

		addrs, err := i.Addrs()
		if err != nil {
			// If we get an error, keep track of it and continue with the next interface
			errs = append(errs, err)
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ipList = append(ipList, v.IP.String())
				if(v.IP.To4()!=nil&&!v.IP.IsLoopback()) {
					if "linux" == sysType {
						maskList = append(maskList, v.IP.String()+"--"+ipv4MaskString(v.Mask)+"--"+hw+"--"+getaway)
					}
					if "windows" != sysType && "linux" != sysType {
						maskList = append(maskList, v.IP.String()+"--"+ipv4MaskString(v.Mask)+"--"+hw)
					}
				}
			case *net.IPAddr:
				ipList = append(ipList, v.IP.String())
			}
		}


	}
	if "windows" == sysType {
		maskList,err=GetMaskList()
	}

	return ipList, hwList,maskList, errs.Err()
}

func ipv4MaskString(m []byte) string {
	if len(m) != 4 {
		fmt.Println("ipv4Mask: len must be 4 bytes")
		return ""
	}

	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}
