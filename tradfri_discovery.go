package tradfri

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/mdns"
)

type DiscoveryInfo struct {
	Name string
	Host string

	Address string
}

func Discover() ([]string, error) {
	addresses := make([]string, 0)

	entries := make(chan *mdns.ServiceEntry, 0)

	go func() {
		for entry := range entries {
			if !strings.HasPrefix(entry.Host, "TRADFRI-Gateway-") || entry.AddrV4 == nil {
				continue
			}

			address := fmt.Sprintf("%s:%d", entry.AddrV4, entry.Port)
			addresses = append(addresses, address)
		}
	}()

	params := &mdns.QueryParam{
		Service: "_coap._udp",
		Domain:  "local",
		Timeout: 1 * time.Second,
		Entries: entries,
	}

	err := mdns.Query(params)

	if err != nil {
		return nil, err
	}

	close(entries)

	return addresses, nil
}
