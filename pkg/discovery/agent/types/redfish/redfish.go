// Copyright 2023 NJWS Inc.

package redfish

import (
	"github.com/koron/go-ssdp"
	"github.com/vishvananda/netlink"
)

const (
	redfishFilter = "urn:dmtf-org:service:redfish-rest:1"
)

func GetServices() (servicesMap map[string]ssdp.Service, err error) {
	servicesMap = make(map[string]ssdp.Service)

	links, err := netlink.LinkList()
	if err != nil {
		return
	}

	for _, link := range links {
		addrs, err := netlink.AddrList(link, 2)
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			services, err := ssdp.Search(redfishFilter, 1, addr.IP.String()+":1900")
			if err != nil {
				continue
			}

			for _, service := range services {
				if service.Type != redfishFilter {
					continue
				}
				servicesMap[service.USN] = service
			}
		}
	}
	return
}
