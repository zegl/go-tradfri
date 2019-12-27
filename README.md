# go-tradfri

Go SDK for IKEA TRÅDFRI Gateway using mDNS, DTLS and CoAP

## Features

- discover gateways using mDNS
- exchange pre-shared-key using security code
- gateway info
- list devices, get device details, update light & outlet settings
- list groups, get group details
- list scenes, get scene details

## Libraries

- https://github.com/pion/dtls
- https://github.com/dustin/go-coap
- https://github.com/hashicorp/mdns

## References

- https://github.com/glenndehaan/ikea-tradfri-coap-docs
- https://github.com/AlCalzone/node-tradfri-client
- https://github.com/ggravlingen/pytradfri

## Examples

discover gatways in network

```go
func main() {
	addresses, err := tradfri.Discover()

	for _, address := range addresses {
		log.WithFields(log.Fields{
			"address": address,
		}).Info("discovered")
    }
}
```

get and store pre-shared-key

```go
func main() {
    psk, err := psk(address, clientID, securityCode)
    
    client, err := tradfri.New(address, clientID, psk)    
    defer client.Close()
}

func psk(address, clientID, securityCode string) ([]byte, error) {
	filename := ".tradfri_" + clientID

	if data, err := ioutil.ReadFile(filename); err == nil {
		return data, nil
	}

	psk, err := tradfri.PSK(address, clientID, securityCode)

	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(filename, []byte(psk), 0644); err != nil {
		return nil, err
	}

	return psk, nil
}
```

set color of a specific bulb

```go
func main() {
    client, err := tradfri.New(address, clientID, psk)    
    defer client.Close()

	color := "dc4b31"

    client.UpdateDevice(deviceID, tradfri.DeviceSettings{
	    LightSettings: []tradfri.LightSettings{
            tradfri.LightSettings{
                Color: &color,
            },
        },
    })
}
```

turn off specific outlet

```go
func main() {
	client, err := tradfri.New(address, clientID, psk)    
	defer client.Close()

	power := 0
	
	client.UpdateDevice(deviceID, tradfri.DeviceSettings{
	 	OutletSettings: []tradfri.OutletSettings{
	 		tradfri.OutletSettings{
	 			Power: &power,
	 		},
	 	},
	})
```

