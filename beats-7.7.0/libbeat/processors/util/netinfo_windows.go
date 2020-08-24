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
// +build windows

package util

import (
	"fmt"
	"github.com/joeshaw/multierror"
	"net"
	"os"
	"syscall"
	"unsafe"
)

func Getaway(){

}
// GetNetInfo returns lists of IPs and MACs for the machine it is executed on.
func GetMaskList() (maskList []string, err error) {
	// Get all interfaces and loop through them
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Keep track of all errors
	var errs multierror.Errors

	for _, i := range ifaces {
		// Skip loopback interfaces
		if i.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		hw := i.HardwareAddr.String()
		aList, err := getAdapterList()
		if err != nil {
			continue
		}
		for ai := aList; ai != nil; ai = ai.Next {
			index := ai.Index

			if i.Index == int(index) {
				ipl := &ai.IpAddressList
				gwl := &ai.GatewayList
				for ; ipl != nil; ipl = ipl.Next  {
					ip :=fmt.Sprintf("%s",ipl.IpAddress.String)
					mask :=fmt.Sprintf("%s",ipl.IpMask.String)
					geteway :=fmt.Sprintf("%s", gwl.IpAddress.String)

					if "0.0.0.0"!=ip {
						maskList=append(maskList,fmt.Sprintf("%s",ip+"--"+mask+"--"+hw+"--"+geteway))
					}

				    fmt.Printf("%s: %s %s %s\n", i.Name, ipl.IpAddress.String, ipl.IpMask.String, gwl.IpAddress.String)
				}
			}
		}

	}

	return maskList, errs.Err()
}

func getAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0]))
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}

