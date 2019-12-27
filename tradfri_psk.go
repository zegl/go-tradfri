package tradfri

import (
	"encoding/json"
	"errors"
	"net"
	"strings"
	"time"

	"github.com/dustin/go-coap"
	"github.com/pion/dtls/v2"
)

func PSK(address, id, code string) ([]byte, error) {
	if code == "" {
		return nil, errors.New("invalid security code")
	}

	if !strings.Contains(address, ":") {
		address = address + ":5684"
	}

	addr, err := net.ResolveUDPAddr("udp", address)

	if err != nil {
		return nil, err
	}

	config := &dtls.Config{
		PSKIdentityHint: []byte("Client_identity"),
		PSK: func(hint []byte) ([]byte, error) {
			return []byte(code), nil
		},
		ConnectTimeout: dtls.ConnectTimeoutOption(30 * time.Second),
		CipherSuites:   []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_CCM_8},
	}

	client, err := dtls.Dial("udp", addr, config)

	if err != nil {
		return nil, err
	}

	t := &Tradfri{
		client: client,
	}

	defer t.Close()

	type requestType struct {
		ID string `json:"9090"`
	}

	type responseTyype struct {
		PSK string `json:"9091"`
	}

	payload, err := json.Marshal(requestType{
		ID: id,
	})

	request := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.POST,
		MessageID: t.MessageID(),
		Payload:   payload,
	}

	request.SetPathString("/15011/9063")

	response, err := t.RoundTrip(request)

	if err != nil {
		return nil, err
	}

	if response.Code != coap.Created {
		return nil, errors.New("response is not of type created")
	}

	var data responseTyype

	if err := json.Unmarshal(response.Payload, &data); err != nil {
		return nil, err
	}

	if len(data.PSK) == 0 {
		return nil, errors.New("unexpected response")
	}

	return []byte(data.PSK), nil
}
