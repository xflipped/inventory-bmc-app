package ssdp

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/koron/go-ssdp"
)

const (
	redfishFilter = "urn:dmtf-org:service:redfish-rest:1"
)

type RedfishService struct {
	Name     string
	URL      string
	Username string
	Password string
}

func Discover() ([]RedfishService, error) {
	redfishServices := []RedfishService{}

	srvs, err := discoverRedfishServices()
	if err != nil {
		return nil, err
	}

	for _, url := range srvs {
		service := RedfishService{
			Name:     uuid.NewString(),
			URL:      url,
			Username: "root",
			Password: "changeme",
		}
		redfishServices = append(redfishServices, service)
	}

	fmt.Printf("Services: %+v\n", redfishServices)

	return redfishServices, nil
}

func discoverRedfishServices() (map[string]string, error) {
	m := make(map[string]string)

	mifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifi := range mifs {

		fmt.Printf("\nInterface name: %s\n", ifi.Name)
		addr, err := getInterfaceIpAddr(ifi.Name)
		if err != nil {
			continue
		}

		fmt.Printf("\nInterface IP address: %s\n", addr)

		if len(addr) > 0 {
			list, err := ssdp.Search(redfishFilter, 1, addr+":1900")
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, srv := range list {
				if srv.Type == redfishFilter {
					u, err := url.Parse(srv.Location)
					if err != nil {
						fmt.Println(err)
						continue
					}
					m[srv.Location] = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				}
			}
		}
	}

	return m, nil
}

func getInterfaceIpAddr(interfaceName string) (addr string, err error) {
	var (
		ief   *net.Interface
		addrs []net.Addr
	)

	if ief, err = net.InterfaceByName(interfaceName); err != nil {
		return
	}

	if addrs, err = ief.Addrs(); err != nil {
		return
	}

	for _, addr := range addrs {
		if ip4, ok := addr.(*net.IPNet); ok {
			if ipv4Addr := ip4.IP.To4(); ipv4Addr != nil {
				return ipv4Addr.String(), err
			}
		}
	}

	err = errors.New(fmt.Sprintf("Interface %s doesn't have an ipv4 address\n", interfaceName))
	return
}
